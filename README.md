# Go Expert - Client Sever API

Este desafio compreende aplicar os conhecimentos adquiridos em: 1. Servidor _HTTP_; 2. Contextos; 3. Bancos de Dados; 4. Manipulação de Arquivos.

Antes, uma pequena nota sobre Contextos.

Contexto é um um recurso do _Go_ que permite controlar o que a aplicação está executando em determinado instante. Ele carrega, por exemplo, limites de tempo (_deadlines_), sinais de cancelamento e outros valores em escopo de _request_. Por exemplo, caso passe do tempo limite (_deadline_) parametrizado para executar um _request_, o Contexto pára a execução da operação.

Outro ponto importante é que o Contexto permite guardar informações para serem resgatadas em outras áreas da aplicação. Em _headers_ _HTTP_, chamadas de filas, etc., é possível obter informações a partir do Contexto, como o _correlation ID_, para trabalhar com _tracing_ distribuído, por exemplo.

Quando a aplicação faz uma chamada _HTTP_ para determinada _API_, caso esteja muito lenta, de forma a não travar o processo da aplicação, o Contexto cancela a operação. Da mesma forma com o banco de dados: caso uma consulta esteja lenta, o Contexto cancela sua execução.

O uso de Contexto é incentivado como boa prática pela própria equipe do _Google_: "No _Google_, exigimos que os programadores _Go_ passem um parâmetro _Context_ como primeiro argumento no caminho da chamada entre _requests_ recebidos e enviados. (...) Isso permite um controle simples de _timeouts_ e cancelamentos e garante que valores críticos (...) transitem de forma correta pelos programas _Go_." (https://go.dev/blog/context).

Este desafio consiste em implementar uma aplicação cliente-servidor:

- O cliente realiza uma requisição _HTTP_ para o servidor solicitando a cotação do Dólar.

- O servidor consume a _API_ contendo o câmbio do Dólar e Real no endereço: `https://economia.awesomeapi.com.br/json/last/USD-BRL` e, em seguida, retorna no formato _JSON_ o resultado para o cliente.

- Usando o _package context_, o servidor registrar no banco de dados _SQLite_ cada cotação recebida. O _timeout_ máximo para a chamada da _API_ de cotação do Dólar deve ser 200ms e o _timeout_ máximo para persistir os dados no banco deve ser de 10ms.

- O cliente recebe do servidor apenas o valor atual do câmbio. Utilizando o _package context_, o cliente tem um _timeout_ máximo de 300ms para receber o resultado do servidor.

- O cliente salva a cotação em um arquivo "cotacao.txt".

- O _endpoint_ gerado pelo servidor deve ser: `/cotacao` e servir na porta 8080.

### Solução

1. Abrir o terminal na raiz do projeto;

2. Subir o servidor, executando:

```
go run server.go
```

3. Ao subir o servidor, é gerado, no diretório raiz, o arquivo **test.db** do SQLite. A cada vez que se sobe o servidor, o banco de dados é iniciado com a tabela **cotacoes** vazia;

4. Abrir um novo terminal e mover-se para o diretório **client**;

5. Executar o comando:

```
go run client.go
```

6. É gerado o arquivo **cotacao.txt** no diretório **client**;

7. Pode-se conferir o registro da cotação recebida no diretório raiz:

```sh
sqlite3 test.db
select * from cotacoes;
```

