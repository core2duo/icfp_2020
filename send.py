import requests

params = {'apiKey': "6d65082372354d349977e0f48cd1c95f"}
url = "https://icfpc2020-api.testkontur.ru/aliens/send"


def send(data):
    r = requests.post(url, params=params, data=data)
    return r.json()