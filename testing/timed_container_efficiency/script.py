import time
import numpy as np
import requests

def send_ping(ip, event):
    timestamp = time.time()
    try:
        requests.post(f"http://{ip}:5000/ping", json={"event": event, "timestamp": timestamp})
        print(f"Sent {event} ping to {ip} at {timestamp}")
    except Exception as e:
        print(f"Failed to send ping: {e}")

def cpu_intensive_task(duration):
    """Perform dummy NumPy CPU operations to simulate work."""
    start_time = time.time()
    while time.time() - start_time < duration:
        _ = np.random.rand(1000, 1000).dot(np.random.rand(1000, 1000))  # Heavy matrix multiplication

if __name__ == "__main__":
    TARGET_IP = "XXXXX"  # Replace with IP of computer you are using to test
    COMPUTE_DURATION = 5  # Duration in seconds

    send_ping(TARGET_IP, "start")
    cpu_intensive_task(COMPUTE_DURATION)
    send_ping(TARGET_IP, "end")