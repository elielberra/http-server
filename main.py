import socket

host = "localhost"
port = 8000

custom_request = (
"""POST /api/users HTTP/1.1
Host: example.com
Content-Type: application/json
Content-Length: 55

{
  "name": "Jane Doe",
  "email": "jane@example.com"
}
"""
)

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.connect((host, port))
    s.sendall(custom_request.encode())
    response = s.recv(4096)
    print(response.decode())