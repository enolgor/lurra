import createClient from "openapi-fetch";
import type { components, paths } from "./schema";

type Todo = components["schemas"]["model.Todo"];

const client = createClient<paths>({ baseUrl: "/api/v1" });

const form = document.querySelector<HTMLFormElement>("#add-form")!;
const titleInput = document.querySelector<HTMLInputElement>("#add-title")!;
const errorEl = document.querySelector<HTMLParagraphElement>("#error")!;
const list = document.querySelector<HTMLUListElement>("#todos")!;

function showError(message?: string) {
  errorEl.textContent = message ?? "";
  errorEl.hidden = !message;
}

function renderTodo(todo: Todo): HTMLLIElement {
  const checkbox = document.createElement("input");
  checkbox.type = "checkbox";
  checkbox.checked = todo.done;
  checkbox.addEventListener("change", () => toggleTodo(todo.id, checkbox.checked));

  const title = document.createElement("span");
  title.textContent = todo.title;

  const label = document.createElement("label");
  label.append(checkbox, title);

  const remove = document.createElement("button");
  remove.textContent = "Delete";
  remove.addEventListener("click", () => deleteTodo(todo.id));

  const item = document.createElement("li");
  item.classList.toggle("done", todo.done);
  item.append(label, remove);
  return item;
}

async function loadTodos() {
  const { data, error } = await client.GET("/todos");
  if (error) return showError("Could not load todos");
  list.replaceChildren(...data.map(renderTodo));
}

async function addTodo(title: string) {
  const { error } = await client.POST("/todos", { body: { title } });
  if (error) return showError(error.message);
  showError();
  await loadTodos();
}

async function toggleTodo(id: number, done: boolean) {
  const { error } = await client.PATCH("/todos/{id}", {
    params: { path: { id } },
    body: { done },
  });
  if (error) showError(error.message);
  else showError();
  await loadTodos();
}

async function deleteTodo(id: number) {
  const { error } = await client.DELETE("/todos/{id}", {
    params: { path: { id } },
  });
  if (error) showError(error.message);
  else showError();
  await loadTodos();
}

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  await addTodo(titleInput.value);
  form.reset();
  titleInput.focus();
});

loadTodos();
