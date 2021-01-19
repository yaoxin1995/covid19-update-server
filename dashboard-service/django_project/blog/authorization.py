import http.client
import ast
# conn = http.client.HTTPSConnection("scc2020g8.eu.auth0.com")

# payload = "{\"client_id\":\"9QDP784vncoXslIJ5H0pFWuQcvxySxxx\",\"client_secret\":\"ANA3qovFC2UfdOkdzxtowaAXiO_oPBO1RCJelEXsy6WJjUwJQSVpr3mPMNM9JcBi\",\"audience\":\"https://185.128.119.135\",\"grant_type\":\"client_credentials\"}"

# headers = { 'content-type': "application/json" }

# conn.request("POST", "/oauth/token", payload, headers)

# res = conn.getresponse()
# data = res.read()

#key={"access_token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik5ZTUxXWGJHV3VOMjJnOGx4djB0QiJ9.eyJpc3MiOiJodHRwczovL3NjYzIwMjBnOC5ldS5hdXRoMC5jb20vIiwic3ViIjoiOVFEUDc4NHZuY29Yc2xJSjVIMHBGV3VRY3Z4eVN4eHhAY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vMTg1LjEyOC4xMTkuMTM1IiwiaWF0IjoxNjEwOTk2OTMwLCJleHAiOjE2MTEwODMzMzAsImF6cCI6IjlRRFA3ODR2bmNvWHNsSUo1SDBwRld1UWN2eHlTeHh4IiwiZ3R5IjoiY2xpZW50LWNyZWRlbnRpYWxzIn0.NTzPVRFsdPd6fro3hvGH7Q7x99O36xZ1HlBWgV-Vw0aLmUzvxqPKsOhnYWmJ6Iu1-OXc8ie2RR9l-P-bwH6qi1DDYQ8nfQBE39QNfuCJnUMeLvCe2o9iB2DNNjHmw8N7WeSas2RUotuk7wI29ghMdVyzZWvEdQbecIMzXqj0n7FclQ9xj0EK4WU4jJDou4vnPWcpUjfto1A-xFz8UBB4j79j4icTqRw45JhhGDsRDnoHDJIcLzduBzDL1kGJe3yQsaWbIRCT3mun__C0RRJn7lZC_pAKEh35bYpRRvR7ERuGUWGl1hNh4mJFoPBqOfE2EIoW7nzlLrZ4Ktmlp9x02Q","expires_in":86400,"token_type":"Bearer"}


#print(data.decode("utf-8"))

def getAuthorization():
    conn = http.client.HTTPSConnection("scc2020g8.eu.auth0.com")
    payload = "{\"client_id\":\"9QDP784vncoXslIJ5H0pFWuQcvxySxxx\",\"client_secret\":\"ANA3qovFC2UfdOkdzxtowaAXiO_oPBO1RCJelEXsy6WJjUwJQSVpr3mPMNM9JcBi\",\"audience\":\"https://185.128.119.135\",\"grant_type\":\"client_credentials\"}"
    headers = { 'content-type': "application/json" ,"accept": "application/json"}
    conn.request("POST", "/oauth/token", payload, headers)
    res = conn.getresponse()
    data = res.read()
    key = data.decode("utf-8")
    key_dic = ast.literal_eval(key)
    return key_dic["access_token"]
   # type : string   key_dic = ast.literal_eval(key)


    #conn1 = http.client.HTTPSConnection("scc2020g8.eu.auth0.com")
