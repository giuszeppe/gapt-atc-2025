import requests
import threading
import time
import matplotlib.pyplot as plt

AUTH_URL = "http://localhost:8080/login"
SIMULATION_URL = "http://localhost:8080/post-simulation"
USERNAME = "admin"
PASSWORD = "password"

CONCURRENT_USERS_LIST = [150, 200, 250, 300, 400, 500, 600, 800, 1000]

def authenticate():
	payload = {"username": USERNAME, "password": PASSWORD}
	resp = requests.post(AUTH_URL, json=payload)
	if resp.status_code != 200:
		raise Exception(f"auth failed: {resp.text}")
	data = resp.json()
	return data["data"]["token"]

def simulate(token, user_id, result, durations, lock):
	payload = {
		"scenario_id": 1,
		"input_type": "text",
		"scenario_type": "takeoff",
		"role": "tower",
		"simulation_advancement_type": "steps",
		"mode": "single"
	}
	headers = {
		"Authorization": token,
		"Content-Type": "application/json"
	}
	start = time.time()
	try:
		resp = requests.post(SIMULATION_URL, headers=headers, json=payload, timeout=10)
		elapsed = time.time() - start
		with lock:
			durations.append(elapsed)
		if resp.status_code == 200:
			result["success"] += 1
		else:
			result["fail"] += 1
	except Exception:
		elapsed = time.time() - start
		with lock:
			durations.append(elapsed)
		result["fail"] += 1

def run_stress(concurrent_users):
	try:
		token = authenticate()
	except Exception as e:
		print(f"Failed to authenticate: {e}")
		return 0, concurrent_users, 0.0  # All fail

	result = {"success": 0, "fail": 0}
	durations = []
	lock = threading.Lock()
	threads = []
	for i in range(concurrent_users):
		t = threading.Thread(target=simulate, args=(token, i+1, result, durations, lock))
		threads.append(t)
		t.start()
	for t in threads:
		t.join()
	avg_time = sum(durations) / len(durations) if durations else 0.0
	return result["success"], result["fail"], avg_time

def main():
	successes = []
	failures = []
	avg_times = []
	for users in CONCURRENT_USERS_LIST:
		print(f"Testing with {users} concurrent users...")
		success, fail, avg_time = run_stress(users)
		successes.append(success)
		failures.append(fail)
		avg_times.append(avg_time)
		print(f"Success: {success}, Fail: {fail}, Avg Time: {avg_time:.3f}s")

	# Plotting
	plt.figure(figsize=(12, 8))
	plt.subplot(2, 1, 1)
	plt.plot(CONCURRENT_USERS_LIST, successes, label="Success", marker='o')
	plt.plot(CONCURRENT_USERS_LIST, failures, label="Fail", marker='x')
	plt.xlabel("Concurrent Users")
	plt.ylabel("Number of Requests")
	plt.title("Stress Test: Success vs Failures by Concurrent Users")
	plt.legend()
	plt.grid(True)

	plt.subplot(2, 1, 2)
	plt.plot(CONCURRENT_USERS_LIST, avg_times, label="Avg Time per Request (s)", marker='s', color='purple')
	plt.xlabel("Concurrent Users")
	plt.ylabel("Avg Time (seconds)")
	plt.title("Average Time per Request")
	plt.legend()
	plt.grid(True)

	plt.tight_layout()
	plt.show()

if __name__ == "__main__":
	main()