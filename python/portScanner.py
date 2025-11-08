import socket
import asyncio
import argparse
import ipaddress
import logging
import json
import sys

# Setup logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def parse_arguments():
    parser = argparse.ArgumentParser(description='Port Scanner for TCP and UDP ports.')
    parser.add_argument('--target', required=True, help='Target IP, CIDR, or start-end range (e.g., 192.168.1.1, 192.168.1.0/24, 192.168.1.1-192.168.1.10)')
    parser.add_argument('--protocols', default='tcp,udp', help='Comma-separated list of protocols to scan (tcp, udp). Default: tcp,udp')
    parser.add_argument('--timeout', type=float, default=1.0, help='Timeout in seconds for each port scan. Default: 1.0')
    parser.add_argument('--concurrency', type=int, default=100, help='Maximum concurrent tasks. Default: 100')
    parser.add_argument('--output', help='Output file to save results in JSON format.')
    return parser.parse_args()

def expand_ip_range(target):
    ips = []
    try:
        if '/' in target:
            # CIDR
            network = ipaddress.ip_network(target, strict=False)
            ips = [str(ip) for ip in network.hosts()]
        elif '-' in target:
            # Range
            start, end = target.split('-')
            start_ip = ipaddress.ip_address(start.strip())
            end_ip = ipaddress.ip_address(end.strip())
            if int(start_ip) > int(end_ip):
                raise ValueError("Start IP must be less than or equal to end IP")
            for ip_int in range(int(start_ip), int(end_ip) + 1):
                ips.append(str(ipaddress.ip_address(ip_int)))
        else:
            # Single IP
            ipaddress.ip_address(target)  # Validate
            ips = [target]
    except ValueError as e:
        logging.error(f"Invalid IP range: {target} - {e}")
        sys.exit(1)
    return ips

async def scan_tcp_port(ip, port, timeout, semaphore):
    async with semaphore:
        logging.debug(f"Starting TCP scan for {ip}:{port}")
        try:
            reader, writer = await asyncio.wait_for(asyncio.open_connection(ip, port), timeout=timeout)
            writer.close()
            await writer.wait_closed()
            logging.debug(f"TCP port {port} open on {ip}")
            return port, True
        except (asyncio.TimeoutError, ConnectionRefusedError, OSError) as e:
            logging.debug(f"TCP port {port} closed on {ip}: {e}")
            return port, False

async def scan_udp_port(ip, port, timeout, semaphore):
    async with semaphore:
        logging.debug(f"Starting UDP scan for {ip}:{port}")
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            sock.settimeout(timeout)
            logging.debug(f"UDP sending to {ip}:{port}")
            sock.sendto(b'', (ip, port))
            # Wait for response or timeout
            logging.debug(f"UDP waiting for response on {ip}:{port}")
            data, addr = sock.recvfrom(1024)
            sock.close()
            logging.debug(f"UDP port {port} open on {ip} (received data)")
            return port, True  # Received data, port open
        except socket.timeout:
            # Timeout, assume open (UDP doesn't guarantee response)
            logging.debug(f"UDP port {port} timeout on {ip}, assuming open")
            sock.close()
            return port, True
        except OSError as e:
            logging.debug(f"UDP port {port} closed on {ip}: {e}")
            sock.close()
            return port, False

class PortScanner:
    def __init__(self, ips, protocols, timeout, concurrency):
        self.ips = ips
        self.protocols = protocols
        self.timeout = timeout
        self.concurrency = concurrency
        self.results = {}

    async def scan_ip(self, ip):
        logging.info(f"Starting scan for IP: {ip}")
        open_ports = {'tcp': [], 'udp': []}
        semaphore = asyncio.Semaphore(self.concurrency)

        # Separate TCP and UDP
        tcp_tasks = []
        udp_tasks = []
        if 'tcp' in self.protocols:
            for port in range(1, 1025):  # Scan common ports 1-1024
                tcp_tasks.append(scan_tcp_port(ip, port, self.timeout, semaphore))
        if 'udp' in self.protocols:
            for port in range(1, 1025):  # Scan common ports 1-1024
                udp_tasks.append(scan_udp_port(ip, port, self.timeout, semaphore))

        logging.info(f"Created {len(tcp_tasks)} TCP tasks and {len(udp_tasks)} UDP tasks for IP {ip}")
        logging.info(f"Starting TCP gather for IP {ip}")
        raw_tcp_results = await asyncio.gather(*tcp_tasks, return_exceptions=True) if tcp_tasks else []
        logging.info(f"TCP gather completed for IP {ip}")
        logging.info(f"Starting UDP gather for IP {ip}")
        raw_udp_results = await asyncio.gather(*udp_tasks, return_exceptions=True) if udp_tasks else []
        logging.info(f"UDP gather completed for IP {ip}")

        for result in raw_tcp_results:
            if isinstance(result, BaseException):
                logging.error(f"TCP scan error: {result}")
            else:
                port, is_open = result
                if is_open:
                    open_ports['tcp'].append(port)

        for result in raw_udp_results:
            if isinstance(result, BaseException):
                logging.error(f"UDP scan error: {result}")
            else:
                port, is_open = result
                if is_open:
                    open_ports['udp'].append(port)

        self.results[ip] = open_ports
        logging.info(f"Completed scan for IP: {ip} - Open TCP: {len(open_ports['tcp'])}, UDP: {len(open_ports['udp'])}")

    async def scan_all(self):
        for ip in self.ips:
            await self.scan_ip(ip)

def main():
    args = parse_arguments()
    protocols = [p.strip().lower() for p in args.protocols.split(',')]
    if not all(p in ['tcp', 'udp'] for p in protocols):
        logging.error("Invalid protocols. Must be 'tcp' or 'udp'")
        sys.exit(1)

    ips = expand_ip_range(args.target)
    if len(ips) > 1000:
        logging.warning(f"Large IP range: {len(ips)} IPs. This may take time.")

    scanner = PortScanner(ips, protocols, args.timeout, args.concurrency)
    asyncio.run(scanner.scan_all())

    # Output results
    if args.output:
        with open(args.output, 'w') as f:
            json.dump(scanner.results, f, indent=4)
        logging.info(f"Results saved to {args.output}")
    else:
        print(json.dumps(scanner.results, indent=4))

if __name__ == '__main__':
    main()