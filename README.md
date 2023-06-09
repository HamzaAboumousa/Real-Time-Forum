
![Logo](./template/img/noback.svg)


# Sport Blog

The Sport Blog project is a web application that utilizes WebSocket technology to enable real-time chat functionality. The backend of the application is developed using the Go programming language (Golang), while the frontend is built with JavaScript.




## Installation

To install the "Sport Blog"  web application, use the following command:

```bash
    git clone https://zone01normandie.org/git/haboumou/real-time-forum.git
```
After cloning the repository, create a .pem file for the SSL certificate. Run the following command in the "/real-time-forum" folder:

```bash
    openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout key.pem -out cert.pem
```
Once the certificate is generated, run the blog using the following command in the same folder:
```bash
    go run main.go
```
Now, open your browser and access the blog at:
```
    https://localhost:8000
 ```

## Screenshots

![App Screenshot](/template/screen.jpg)

