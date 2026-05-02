# AGT_SSR — Astro + Go `html/template` SSR (POC)

POC showing **Astro static build** + **Go SSR (`html/template`)** working together.

Repo: https://github.com/Mboukhal/AGT_SSR

---

## 🧠 Concept

```text
Astro → build HTML (keeps {{ .Data }} as-is)
        ↓
Go → parses HTML with html/template → injects data (SSR)
```

---

## 📍 Demo in this repo

* Frontend (Astro page): `./ui/src/pages/login.astro`
* Backend (Go SSR): `./cmd/main.go`

---

## ✏️ Example (from login page)

```astro
{"{{.Name}}"}
<If condition={import.meta.env.DEV}>
  <span class="text-sm text-gray-500">(Development Mode)</span>
  <Else />
  <span class="text-sm text-gray-500">(Production Mode)</span>
</If>
```

## ▶️ Run Demo

```bash
# 1. clone
git clone https://github.com/Mboukhal/AGT_SSR
cd AGT_SSR

# 2. install using bun
make install
mv .env.exemple .env

# 3. run Project
make 
```

Open:

```
http://localhost:3000
```


## 👤 Credits

Demo by @Mboukhal


### What happens

* `{"{{.Name}}"}` → Astro outputs literal `{{ .Name }}` into HTML
* `<If ...>` → resolved at **build time** by Astro
* Final HTML still contains Go template syntax

---

## 🧩 Go SSR (simplified)

```go
tpl := template.Must(template.ParseFiles("ui/dist/login/index.html"))

tpl.Execute(w, map[string]any{
	"Name": "Mohammed",
})
```

---

## ⚠️ Important rules

* Astro must **NOT process Go templates**
* Use `{"{{...}}"}` when inside `.astro`
* Disable HTML minify OR use a template-safe minifier
* Go handles all runtime rendering

---

## 🚀 Result

* Astro = UI + components (React/Svelte/Vue/MDX)
* Go = SSR engine
* Clean separation, fast runtime

---

## 🧠 Summary

```text
Astro = build-time
Go = runtime
```

---

## 🔥 Notes

* Works with conditions, loops, variables (`{{ if }}`, `{{ range }}`)
* Safe for SSR with `html/template`
* Dev/Prod logic handled by Astro (`import.meta.env.DEV`)

---

## 📄 License

MIT

@theprimeagen
