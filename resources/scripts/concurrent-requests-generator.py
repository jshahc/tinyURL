import requests
import uuid
from concurrent.futures import ThreadPoolExecutor, as_completed


def generate_random_url():
    random_uuid = uuid.uuid4()
    random_url = f"https://example.com/{random_uuid}"
    return random_url


def make_post_request():
    data = {'url': generate_random_url()}
    response = requests.post('http://localhost:8080/', json=data)
    return response.text


def main():
    # Using ThreadPoolExecutor to make parallel POST requests
    with ThreadPoolExecutor(max_workers=50) as executor:
        # Submitting tasks for each URL
        futures = [executor.submit(make_post_request) for i in range(5000)]

        # Collecting results as they become available
        for future in as_completed(futures):
            result = future.result()
            print(result)


if __name__ == "__main__":
    main()
