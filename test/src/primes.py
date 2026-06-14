#!/usr/bin/env python3
import math

def is_prime(n):
    """Returns True if n is prime, False otherwise."""
    if n <= 1:
        return False
    if n == 2:
        return True
    if n % 2 == 0:
        return False
        
    # Check odd factors up to the square root of n
    for i in range(3, int(math.sqrt(n)) + 1, 2):
        if n % i == 0:
            return False
    return True

def display_primes(limit):
    """Finds and prints all prime numbers up to the specified limit."""
    print(f"\nPrime numbers up to {limit}:")
    
    count = 0
    for num in range(2, limit + 1):
        if is_prime(num):
            # Print numbers side-by-side separated by spaces
            print(num, end=" ")
            count += 1
            
    print(f"\n\nTotal prime numbers found: {count}")

if __name__ == "__main__":
    try:
        user_limit = int(input("Enter the upper limit to find prime numbers: "))
        display_primes(user_limit)
    except ValueError:
        print("Please enter a valid integer.")

