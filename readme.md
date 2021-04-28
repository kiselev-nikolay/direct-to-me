# Direct to me

Self-hosted application for redirecting user or application data in the background. Source code written in Go, TypeScript. React frontend.

No, it's a data redirection program. You create a link inside the program `example.com/test-link` that will take data in any format and send it in the background to some web hook (for example send it to your slack), and redirect the user to other links, such as back to the home page.

👨‍🏭 Looking for contributors!

### Run in docker

```shell
docker pull ghcr.io/kiselev-nikolay/direct-to-me/direct-to-me:latest
docker run --name direct-to-me -p 8080:8080 ghcr.io/kiselev-nikolay/direct-to-me/direct-to-me:latest
```

Now check this out! Go to [localhost:8080](http://localhost:8080)

## Screenshots 😍

![scr1.png](./docs/scr1.png)

![scr2.png](./docs/scr2.png)

![scr3.png](./docs/scr3.png)

![scr4.png](./docs/scr4.png)

![scr5.png](./docs/scr5.png)

![scr6.png](./docs/scr6.png)
