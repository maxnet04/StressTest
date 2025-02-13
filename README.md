# stress-test

Sistema CLI em Go para realizar testes de carga em um serviço web. O usuário deverá fornecer a URL do serviço, o número total de requests e a quantidade de chamadas simultâneas.

### Como Utilizar localmente:
#### Requisitos:
    - Certifique-se de ter o Go instalado em sua máquina.
    - Certifique-se de ter o Docker instalado em sua máquina.
    
- [GO](https://golang.org/doc/insttall) 1.17 ou superior
- [Docker](https://docs.docker.com/get-docker/)


Como Rodar localmente

  1. Clonar o Repositório:~
  ```git clone https://github.com/maxnet04/StressTest.git```


  2. Acesse a pasta do app:
  ```cd StressTest```

  3. Criar a imagem docker:
  ```docker build -t stress-test-commandd . ```


#### Para testar execute o comando abaixo:

```docker run stress-test-command --url=http://google.com --requests=10 --concurrency=5 ```

## Resultado agrupado por status code:

Performing load test on http://google.com with 10 requests and 5 concurrent cells.
==============================================
Total execution time: 1.396605461s
Total number of request made: 10

Summary:
Status code 200 | Count  10 (100.00%)