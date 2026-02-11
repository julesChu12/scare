
import redis
import sys

def get_sms_code(phone, password):
    try:
        r = redis.Redis(host='localhost', port=6379, password=password, decode_responses=True)
        # Try to find keys matching the phone number
        keys = r.keys(f"*{phone}*")
        print(f"Keys found for {phone}: {keys}")
        
        for key in keys:
            val = r.get(key)
            print(f"Value for {key}: {val}")
            
    except Exception as e:
        print(f"Error connecting to Redis: {e}")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python get_redis_code.py <phone>")
        sys.exit(1)
    
    phone = sys.argv[1]
    password = "scare_redis_pass"
    get_sms_code(phone, password)
