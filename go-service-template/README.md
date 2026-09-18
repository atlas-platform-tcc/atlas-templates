# Serviço Atlas (Go)

Serviço HTTP em Go gerado pelo **Atlas** (Golden Path). Este README é parte do
template — a plataforma exige que todo serviço nasça documentado (conformidade
`template-readme`, ADR atlas-api/0013).

## Endpoints
- `GET /healthz` — *readiness/liveness* (usado pelas probes do chart base).
- `GET /dbz` — alcançabilidade do banco, quando `database.enabled` (a plataforma injeta a conexão).
- `GET /` — resposta de exemplo.

## Como rodar localmente
```bash
go run .
# PORT=8080 por padrão; DB_HOST/DB_PORT + POSTGRES_* injetados pela plataforma quando há banco.
```

## Ciclo de vida (não editar à mão o que a plataforma gera)
- **CI (`.github/workflows/ci.yml`)**: a cada push em `main` → lint → test → build → scan → publica
  a **imagem imutável** (`tag = SHA`). Não faz deploy.
- **Deploy day-2 (`.github/workflows/deploy.yml`)**: botão *Run workflow* promove uma imagem já
  construída (fixa o SHA no `atlas-gitops`; Argo CD reconcilia). É o único gate humano do day-2.
- **Configuração de runtime**: edite apenas `deploy/values.yaml` (réplicas, env). Os guardrails
  (probes, requests/limits, não-root, labels) vêm do chart base e não são configuráveis aqui.

> Gerado pelo Atlas. Ajuste o código da aplicação à vontade; a esteira e os padrões são impostos
> pela plataforma.
