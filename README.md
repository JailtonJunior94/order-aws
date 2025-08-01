# Order AWS

Um projeto Go para gerenciamento de pedidos utilizando serviços AWS com suporte a desenvolvimento local através do LocalStack.

## 📋 Pré-requisitos

Antes de começar, certifique-se de ter instalado:

- [Go](https://golang.org/dl/) 1.21 ou superior
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Terraform](https://www.terraform.io/downloads.html) 0.13.1 ou superior
- [Make](https://www.gnu.org/software/make/)

## 🏗️ Arquitetura

O projeto utiliza os seguintes serviços AWS:

- **DynamoDB**: Armazenamento de dados de pedidos
- **S3**: Bucket para armazenamento de arquivos relacionados aos pedidos
- **SQS**: Fila para processamento assíncrono de pedidos (com Dead Letter Queue)
- **VPC**: Rede virtual para isolamento dos recursos

## 🚀 Configuração do Ambiente

### 1. Configuração de Variáveis de Ambiente

```bash
make dotenv
```

Este comando cria um arquivo `.env` a partir do template `.env.example` no diretório `cmd/`.

### 2. Inicialização do LocalStack

Para iniciar o ambiente local com LocalStack:

```bash
make start_localstack
```

Este comando:
- Inicia o container LocalStack na porta 4566
- Configura a região padrão como `us-east-1`
- Cria um volume persistente para os dados

### 3. Criação da Infraestrutura Local

Para provisionar a infraestrutura AWS local:

```bash
make create_infra_local
```

Este comando:
- Inicializa o Terraform
- Aplica automaticamente a configuração da infraestrutura
- Cria os recursos: DynamoDB, S3, SQS, VPC e Security Groups

### 4. Compilação da Aplicação

Para compilar a aplicação:

```bash
make build_order
```

O binário será gerado em `./bin/order`.

## 🧪 Testes e Qualidade do Código

### Executar Testes

```bash
make test
```

Executa todos os testes com:
- Verificação de condições de corrida (`-race`)
- Cobertura de código
- Perfil de cobertura salvo em `coverage.out`

### Relatório de Cobertura

```bash
make cover
```

Gera um relatório HTML da cobertura de código.

### Análise de Código (Linting)

```bash
make lint
```

Executa o `golangci-lint` com a configuração definida em `.golangci.yml`.

### Verificação de Vulnerabilidades

```bash
make vulncheck
```

Executa o `govulncheck` para identificar vulnerabilidades conhecidas nas dependências.

### Geração de Mocks

```bash
make mocks
```

Gera mocks usando o Mockery para testes unitários.

## 🛠️ Comandos Makefile Disponíveis

| Comando | Descrição |
|---------|-----------|
| `start_localstack` | Inicia o LocalStack em container Docker |
| `stop_localstack` | Para o container LocalStack |
| `create_infra_local` | Cria a infraestrutura AWS local com Terraform |
| `destroy_infra_local` | Destrói a infraestrutura AWS local |
| `dotenv` | Cria arquivo .env a partir do template |
| `build_order` | Compila a aplicação Go |
| `lint` | Executa análise de código |
| `mocks` | Gera mocks para testes |
| `test` | Executa testes unitários |
| `cover` | Gera relatório de cobertura |
| `vulncheck` | Verifica vulnerabilidades |

## 📁 Estrutura do Projeto

```
.
├── Makefile                    # Comandos de automação
├── README.md                   # Documentação do projeto
├── .golangci.yml              # Configuração do linter
├── cmd/                       # Aplicação principal
├── config/                    # Configurações
├── deployment/                # Arquivos de deploy
│   ├── docker-compose.yml     # LocalStack setup
│   └── terraform/             # Infraestrutura como código
│       ├── dynamodb.tf        # Configuração DynamoDB
│       ├── s3.tf              # Configuração S3
│       ├── sqs.tf             # Configuração SQS
│       ├── vpc.tf             # Configuração VPC
│       ├── security_group.tf  # Security Groups
│       ├── variables.tf       # Variáveis Terraform
│       └── providers.tf       # Provedores AWS
├── internal/                  # Código interno da aplicação
└── pkg/                       # Pacotes reutilizáveis
```

## 🔧 Desenvolvimento

### Fluxo de Desenvolvimento Recomendado

1. **Configurar ambiente**:
   ```bash
   make dotenv
   make start_localstack
   make create_infra_local
   ```

2. **Desenvolver e testar**:
   ```bash
   make test
   make lint
   make vulncheck
   ```

3. **Compilar**:
   ```bash
   make build_order
   ```

4. **Limpar ambiente** (quando necessário):
   ```bash
   make destroy_infra_local
   make stop_localstack
   ```

### Recursos da Infraestrutura

- **DynamoDB Table**: `local_orders` com chave hash `id`
- **S3 Bucket**: `local-orders-bucket`
- **SQS Queue**: `local_orders` com DLQ `local_orders_dlq`
- **Endpoints LocalStack**: Disponíveis em `http://localhost:4566`

## 📊 Monitoramento e Logs

O LocalStack fornece logs detalhados dos serviços AWS simulados. Para visualizar:

```bash
docker logs localstack -f
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença [MIT](LICENSE).

## 🆘 Solução de Problemas

### LocalStack não inicia
- Verifique se o Docker está rodando
- Certifique-se de que a porta 4566 não está sendo usada por outro serviço

### Erro no Terraform
- Verifique se o LocalStack está rodando antes de executar `make create_infra_local`
- Execute `make destroy_infra_local` e depois `make create_infra_local` para recriar a infraestrutura

### Testes falhando
- Execute `make mocks` para regenerar os mocks
- Verifique se todas as dependências estão instaladas

## 📞 Suporte

Para dúvidas ou problemas, abra uma issue no repositório do projeto.
