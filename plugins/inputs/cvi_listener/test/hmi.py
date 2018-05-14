from flask import Flask
from flask import request
app = Flask(__name__)

@app.route('/api/v1/results', methods=['PUT'])
def result():
    print(request.data)
    return ''

@app.route('/api/v1/status', methods=['PUT'])
def status():
    print(request.data)
    return ''

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8000)