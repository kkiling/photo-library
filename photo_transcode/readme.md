### Python venv
Create venv
```
python -m venv venv
```
Activate venv
```
.\venv\Scripts\activate
```
Save install packages to requirements
```
pip3 freeze > requirements.txt
```
Install packages from requirements
```
pip3 install -r requirements.txt
```
Pip pakages
```commandline
pip3 install tensorflow
pip3 install keras-applications
pip3 install keras-preprocessing
pip3 install Pillow
pip3 install flask
```

### Run server
```
flask --app server run
```