import socket

host = "localhost"
port = 8000

body = """{
  "name": "Jane Doe",
  "email": "jane@example.com"
}"""

content_length = len(body.encode('utf-8'))

custom_request = f"""POST /api/users HTTP/1.1\r
Host: example.com\r
Content-Type: application/json\r
Content-Length: {content_length}\r
\r
{body}"""

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.connect((host, port))
    s.sendall(custom_request.encode('utf-8'))
    response = s.recv(4096)
    print(response.decode())