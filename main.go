package main

import (
	"html"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var todoListMap = make(map[string][]string)

// sessionIdに対応したtodoListを返す。なければ新規作成して返す
func getTodoList(sessionId string) []string {
	todoList, ok := todoListMap[sessionId]
	if !ok {
		todoList = []string{}
		todoListMap[sessionId] = todoList
	}
	return todoList
}

// /todo ハンドラ
func handleTodo(w http.ResponseWriter, r *http.Request) {
	
	// sessionIdを取得
	sessionId, err := ensureSession(w, r)

	// エラーが発生したら500エラーを返す
	if err != nil{
		http.Error(w, err.Error(), 500)
		return
	}

	// sessionIdに対応したTODOリストを取得
	todoList := getTodoList(sessionId)
	t,_ := template.ParseFiles("templates/todo.html")
	t.Execute(w, todoList)
}

// /add ハンドラ
func handleAddTodo(w http.ResponseWriter, r *http.Request) {

	// sessionIdを取得
	sessionId, err := ensureSession(w, r)
	if err != nil{
		http.Error(w, err.Error(), 500)
		return
	}
	
	// sessionIdに対応したTODOリストを取得
	todoList := getTodoList(sessionId)

	r.ParseForm()
	todo:=  strings.TrimSpace(html.EscapeString(r.Form.Get("todo")))

	if todo != "" {
		todoListMap[sessionId] = append(todoList, todo)
	}

	http.Redirect(w, r, "/todo", 303)
}

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/todo", handleTodo)

	http.HandleFunc("/add", handleAddTodo)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("failed to start:", err)
	}
}