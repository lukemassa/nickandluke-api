#!/usr/bin/env python3

import requests
import argparse

BASE_URL="https://nickandluke-api.herokuapp.com"

def one_test(url, expected_response):
    print(f"Testing {url} ... ", end="")
    actual_response = requests.get(url).json()
    if actual_response == expected_response:
        print("OK")
        return
    print("ERROR")
    print(f"   Expected {expected_response}, got {actual_response}")


def main():
    base_url = BASE_URL

    parser = argparse.ArgumentParser()
    parser.add_argument("--localhost",action="store_true",default=False)
    args = parser.parse_args()

    if args.localhost:
        base_url = "http://localhost:3000"

    one_test(f"{base_url}/guest?name=foobar", {"valid":False,"form":""})
    one_test(f"{base_url}/guest?name=vincent+massa", {'valid': True, 'form': 'https://docs.google.com/forms/d/e/1FAIpQLSevxS_HMScw6Nhcru3ke8GeqWfJnBAA_AdWPc-1eRmgS4G6LQ/viewform?usp=sf_link'})
    one_test(f"{base_url}/guest?name=nancy+massa", {'valid': True, 'form': 'https://docs.google.com/forms/d/e/1FAIpQLSdXF80AevtDqkC7ZTynrzXRuwfZCjQPTpsLhCEfuRPSOCCgww/viewform?usp=sf_link'})

if __name__ == "__main__":
    main()
