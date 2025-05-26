# EditorTexto
Editor de texto em terminal feito com Go!
📝 Projeto: Editor de Texto em Terminal com Go

Este projeto é um editor de texto simples para terminal, desenvolvido em Go, com funcionalidades básicas de inserção, deleção e navegação entre linhas. A estrutura de dados principal utilizada é uma lista duplamente encadeada, que permite manipulação eficiente das linhas do texto.
✨ Funcionalidades

    Inserção de caracteres em qualquer posição

    Navegação por linhas e colunas com as setas do teclado

    Criação e deleção de linhas

    Suporte a comandos como Enter, Backspace, Delete

    Interface responsiva no terminal

📦 Dependências

O projeto utiliza a biblioteca tcell (github.com/gdamore/tcell/v2), uma poderosa abstração de terminal em Go. Essa biblioteca fornece:

    Controle total do terminal com suporte a teclado, mouse e redimensionamento de tela

    Compatibilidade com diferentes ambientes (Unix, Windows, etc.)

    Ricos recursos de renderização com suporte a cores e atributos de estilo

🚀 Execução

Para rodar o editor:

go run main.go
