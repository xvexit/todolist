package main

import (
	"ToDoList/internal/entity"
	database "ToDoList/internal/infrastructure"
	"ToDoList/internal/usecase"
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

func main() {
	ctx := context.Background()
	conn, err := database.InitDatabase(ctx)
	if err != nil {
		panic(err)
	}
	repo := database.NewPostgresRepo(conn)
	uc := usecase.NewUsecase(ctx, repo)
	defer conn.Close(ctx)
	fmt.Println("=== ToDo List Application ===")
	fmt.Println("Введите 'help' для просмотра доступных команд")

	scan := bufio.NewScanner(os.Stdin)

	for scan.Scan() {
		fmt.Print("> ")
		
		command := strings.Fields(scan.Text())

		switch command[0] {
		case "help":
			printHelp()
		case "add":
			if len(command) > 2 {
				if err := uc.AddTask(ctx, command[1], strings.Join(command[2:], " ")); err != nil {
					fmt.Println("Ошибка", err)
				}
			} else {
				printErr()
			}

		case "list":
			if len(command) == 1 {
				rows, err := uc.TaskList(ctx)
				if err != nil {
					fmt.Println("Ошибка", err)
				}
				PrintTaskListPretty(rows)
			} else {
				printErr()
			}
		case "del":
			if len(command) == 2 {
				if err := uc.DelTask(ctx, command[1]); err != nil {
					fmt.Println("Ошибка:", err)
				}
			} else {
				printErr()
			}
		case "done":
			if len(command) == 2 {
				if err := uc.DoTask(ctx, command[1]); err != nil {
					fmt.Println("Ошибка:", err)
				}
			} else {
				printErr()
			}
		case "listpag":
			if len(command) == 2 {
				pageSize, err := strconv.Atoi(command[1])
				if err != nil {
					fmt.Println("Ошибка:", err)
				}
				mapa, err := uc.GetAllTasksInPages(ctx, pageSize)
				if err != nil {
					fmt.Println("Ошибка:", err)
				}
				PrintTaskPagesPretty(mapa)
			} else {
				printErr()
			}
		case "exit":
			return
		default:
			printErr()
		}
	}
}

func printErr() {
	fmt.Println("Неверная команда или количество вводных. Введите help для получения дополнительной информации.")
}

func printHelp() {
	helpText := `
    Доступные команды:
    help    - показать эту справку
    add     - добавить задачу (add {заголовок} {текст})
    list    - список задач
    done    - отметить как выполненное (done {заголовок})
    del     - удалить задачу (del {заголовок})
	listpag - вывод задач постранично (listpag {кол-во элеменитов в 1 странице})
    exit    - выход из программы
    `
	fmt.Println(helpText)
}

func ClearTerminal() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout

	cmd.Run()
}

func printTaskPretty(t entity.Task) {
	status := "❌"
	doneTime := "-"

	if t.Is_done {
		status = "✅"
		if t.Time_done != nil {
			doneTime = t.Time_done.Format("02.01.2006 15:04")
		}
	}

	fmt.Printf(`
%s  %s
────────────────────────────
Описание : %s
Создано  : %s
Выполнено: %s
`,
		status,
		t.Name,
		t.Text,
		t.Time_add.Format("02.01.2006 15:04"),
		doneTime,
	)
}

func PrintTaskListPretty(tasks []entity.Task) {
	if len(tasks) == 0 {
		fmt.Println("📭 Список задач пуст")
		return
	}

	fmt.Println("📋 Список задач")
	fmt.Println("────────────────────────────")

	for i, t := range tasks {
		fmt.Printf("\n[%d]\n", i+1)
		printTaskPretty(t)
	}
}

func PrintTaskPagesPretty(pages map[int][]entity.Task) {
	if len(pages) == 0 {
		fmt.Println("📭 Задачи отсутствуют")
		return
	}

	fmt.Println("📚 Список задач (постранично)")
	fmt.Println("════════════════════════════")

	// чтобы вывод был по порядку
	keys := make([]int, 0, len(pages))
	for k := range pages {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, page := range keys {
		tasks := pages[page]

		fmt.Printf("\n📄 Страница %d\n", page+1)
		fmt.Println("────────────────────────────")

		if len(tasks) == 0 {
			fmt.Println("  (пусто)")
			continue
		}

		for i, t := range tasks {
			fmt.Printf("\n[%d.%d]\n", page+1, i+1)
			printTaskPretty(t)
		}
	}
}
