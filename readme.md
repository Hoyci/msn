# 🛡️ Auth Service - Go + PostgreSQL

Este projeto é um serviço backend escrito em Go com foco em autenticação segura via JWT, gerenciamento de sessões, usuários e categorias (com subcategorias). A arquitetura é modular e preparada para escalabilidade com boas práticas como injeção de dependência, middlewares, tokens RSA e validações robustas.

---

## 🧱 Tecnologias Utilizadas

- **Go 1.21+**
- **PostgreSQL**
- **JWT com chave RSA**
- **Chi Router**
- **SQLX**
- **bcrypt para hashing de senhas**
- **Viper para configuração**
- **Golang Migrate**

---

## 🚀 Funcionalidades

- [x] Registro e autenticação de usuários
- [x] Login com Access/Refresh Token
- [x] Renovação de token via refresh token (com cookies HTTP-only)
- [x] Logout e invalidação de sessões
- [x] Middleware de autenticação
- [x] Categorização de usuários (com ícones e subcategorias)
- [x] Log estruturado com slog (JSON ou modo "bonito" para dev)
- [x] Migrations automáticas via CLI (`make migrate-up`, `make migrate-down`)

---

## 📂 Estrutura do Projeto

```
cmd/                # Entrada da aplicação
internal/
  config/           # Configurações (env, parsing RSA keys)
  infra/            # Infraestrutura: DB, HTTP server, logging
  modules/          # Domínios da aplicação: user, auth, session, categories
pkg/                # Pacotes reutilizáveis: DTOs, utils, faults
```

---

## 🧪 Rodando Localmente

1. Clone o projeto:

```bash
git clone https://github.com/seu-usuario/nome-do-projeto.git
cd nome-do-projeto
```

2. Copie o `.env.example`:

```bash
cp .env.example .env
```

3. Gere suas chaves RSA e adicione ao `.env`:

```bash
# Gere com:
openssl genpkey -algorithm RSA -out access.key
openssl rsa -in access.key -pubout -out access.pub
```

4. Suba o banco com Docker:

```bash
docker-compose up -d postgres
```

5. Rode as migrations:

```bash
make migrate-up
```

6. Inicie o servidor:

```bash
go run cmd/main.go
```

---

## 🔐 Autenticação

- A autenticação utiliza JWT com chave privada (RSA).
- O `access token` é enviado no header `Authorization: Bearer <token>`.
- O `refresh token` é armazenado como cookie `HttpOnly`.

---

## 🗂️ Endpoints

| Método | Rota                   | Descrição                 |
| ------ | ---------------------- | ------------------------- |
| POST   | `/api/v1/auth/login`   | Login do usuário          |
| PATCH  | `/api/v1/auth/logout`  | Logout do usuário         |
| POST   | `/api/v1/auth/refresh` | Renovação de access token |
| GET    | `/api/v1/categories`   | Listar categorias         |

---

## 🧼 Contribuindo

Pull requests são bem-vindos! Para grandes mudanças, abra uma issue primeiro para discutirmos o que você gostaria de mudar.

---

## 📄 Licença

Este projeto está sob a licença MIT.
