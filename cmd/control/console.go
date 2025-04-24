package control

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

const (
	containerdSock  = "/run/containerd/containerd.sock"
	systemNamespace = "system"
	userNamespace   = "user"
	consoleImage    = "docker.io/library/debian:latest"
)


func consoleSwitch(newConsole string) error {
	client, err := containerd.New(containerdSock)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к containerd: %v", err)
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), systemNamespace)

	// Загрузка образа
	image, err := client.Pull(ctx, consoleImage, containerd.WithPullUnpack)
	if err != nil {
		return fmt.Errorf("не удалось загрузить образ: %v", err)
	}

	container, err := client.NewContainer(
		ctx,
		"console-container",
		containerd.WithNewSnapshot("console-snapshot", image),
		containerd.WithNewSpec(
			oci.WithImageConfig(image),
			oci.WithPrivileged, // Запуск контейнера с привилегиями
		),
	)
	if err != nil {
		return fmt.Errorf("не удалось создать контейнер: %v", err)
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return fmt.Errorf("не удалось создать задачу: %v", err)
	}

	if err := task.Start(ctx); err != nil {
		return fmt.Errorf("не удалось запустить задачу: %v", err)
	}

	statusC, err := task.Wait(ctx)
	if err != nil {
		return fmt.Errorf("ошибка ожидания завершения задачи: %v", err)
	}
	<-statusC

	fmt.Printf("Консоль переключена на: %s\n", newConsole)
	return nil
}

func consoleEnable(newConsole string) error {
	client, err := containerd.New(containerdSock)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к containerd: %v", err)
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), systemNamespace)

	image, err := client.Pull(ctx, consoleImage, containerd.WithPullUnpack)
	if err != nil {
		return fmt.Errorf("не удалось загрузить образ: %v", err)
	}

	container, err := client.NewContainer(
		ctx,
		"console-container",
		containerd.WithNewSnapshot("console-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return fmt.Errorf("не удалось создать контейнер: %v", err)
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	fmt.Printf("Консоль %s включена для следующего запуска\n", newConsole)
	return nil
}

func consoleList() error {
	client, err := containerd.New(containerdSock)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к containerd: %v", err)
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), systemNamespace)
	containers, err := client.Containers(ctx)
	if err != nil {
		return fmt.Errorf("не удалось получить список контейнеров: %v", err)
	}

	fmt.Println("Доступные консоли:")
	for _, c := range containers {
		fmt.Printf(" - %s\n", c.ID())
	}

	return nil
}
