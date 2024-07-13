# LoadBalancer

Here's a quick breakdown of how it works:

Load Balancer Initialization: The load balancer is initialized with a list of backend servers (in this case, Facebook, YouTube, and DuckDuckGo).
Round-Robin Load Balancing: Each incoming request is forwarded to one of the backend servers in a round-robin fashion. This ensures even distribution of traffic across all servers.
Reverse Proxy: The reverse proxy component forwards the HTTP requests to the selected backend server, making sure they reach their intended destination.
