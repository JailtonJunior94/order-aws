# Order AWS

Um projeto em Go para gerenciamento de pedidos, utilizando serviços AWS simulados localmente com LocalStack. O objetivo é facilitar o desenvolvimento, testes e integração de aplicações que dependem de recursos AWS, sem custos ou dependência da nuvem real.

## 🚀 Tecnologias Utilizadas

- **Go**: Linguagem principal do projeto.
- **Docker & Docker Compose**: Para orquestração de ambientes e serviços locais.
- **LocalStack**: Simulação de serviços AWS (DynamoDB, S3, SQS, VPC, etc).
- **Terraform**: Infraestrutura como código para provisionamento dos recursos AWS simulados.
- **Taskfile**: Automação de tarefas comuns de desenvolvimento.
- **GolangCI-Lint**: Ferramenta de linting para Go.
- **Mockery**: Geração de mocks para testes.
- **govulncheck**: Verificação de vulnerabilidades em dependências Go.

## 📋 Pré-requisitos

- [Go](https://golang.org/dl/) 1.21 ou superior
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Terraform](https://www.terraform.io/downloads.html) 0.13.1 ou superior
- [Task](https://taskfile.dev/#/installation) (Taskfile runner)

## 🏗️ Arquitetura

O projeto utiliza os seguintes serviços AWS simulados:

- **DynamoDB**: Armazenamento de dados de pedidos
- **S3**: Bucket para arquivos relacionados aos pedidos
- **SQS**: Fila para processamento assíncrono de pedidos (com Dead Letter Queue)
- **VPC**: Rede virtual para isolamento dos recursos

## ⚡ Como Usar

### 1. Configurar variáveis de ambiente

```bash
task dotenv
```
Cria o arquivo `.env` a partir do template `.env.example` no diretório `cmd/`.

### 2. Iniciar o LocalStack

```bash
task start_localstack
```
Inicia o LocalStack em container Docker na porta 4566.

### 3. Provisionar infraestrutura local

```bash
task create_infra_local
```
Inicializa o Terraform e cria os recursos AWS simulados (DynamoDB, S3, SQS, VPC).

### 4. Compilar a aplicação

```bash
task build_order
```
Gera o binário em `./bin/order`.

### 5. Executar testes

```bash
task test
```
Executa todos os testes com verificação de condições de corrida e cobertura de código.

### 6. Gerar relatório de cobertura

```bash
task cover
```
Gera um relatório HTML da cobertura de código.

### 7. Análise de código (lint)

```bash
task lint
```
Executa o `golangci-lint` conforme configuração.

### 8. Verificar vulnerabilidades

```bash
task vulncheck
```
Executa o `govulncheck` para identificar vulnerabilidades conhecidas.

### 9. Gerar mocks para testes

```bash
task mocks
```
Gera mocks usando o Mockery.

### 10. Parar e limpar ambiente

```bash
task stop_localstack
task destroy_infra_local
```
Para o LocalStack e destrói a infraestrutura provisionada.

## 📁 Estrutura do Projeto

```
.
├── cmd/            # Aplicação principal
├── configs/        # Configurações
├── deployment/     # Docker Compose e Terraform
├── internal/       # Domínio e lógica de negócio
├── pkg/            # Pacotes reutilizáveis
├── bin/            # Binários gerados
├── Taskfile.yml    # Automação de tarefas
└── README.md       # Documentação
```

## 🆘 Suporte

Para dúvidas ou problemas, abra uma issue no repositório do projeto.

---

Sinta-se à vontade para adaptar conforme os comandos e detalhes específicos do seu [Taskfile.yml](Taskfile.yml)