## This is a project in go lang for iot backend dashbord
#### Tech used 
  - go lang for writing backend 
  - html for templets
  - css for page style
  - java script for scripts
  - plotly.js for plots in dashbord
  - json files to store creadentials and tresholds
  - csv file to store sensor data
  - smtp for mail comminication
  - fast2sms for message comminication 

#### description of my project
  This is the project we create a api and web dashbord using go lang so that user can view the real time sensor data through charts.
  The authorized persons will be notifed when the values get higher then the treshold value.
  The theshold value will be dynamicaly updated occrding to the users convinence.
  The contact details will be added dynamicaly.
  Uses asyncronus funtions to send mail and sms so thatit doesn't affect the process in our webpage.

## For arduino

#### components used
  - esp 32 dev
  - dht22 sensor
  - sound sensor
  - buzzer
  - LED
  - Breadbord

#### connectons
  - connect the 5v and gnd for dht22 and sound sensor.
  - connect the signal pin of dht22 to gp4 pin, and signal pin of sound sensor to the gp34.
  - now connect the LED to gp25 and buzzer to the gp26 pins and also connect the gnd to them.
  - run the ardino code in the project.


## Set Up and run this project
  - install go lang in your system.
  - clone this repositry into your system.
  - navigate to "notify.go" there enter your email and app password(which will be generated from your mail account).
  - In the same file change the fast2sms key to your key.
  - The run the go project running this command "go run ."
  - Now in find the arduino folder there you  can get the skech for the iot device.
  - On that skech replace the wifi details such as ssid,pasword.
  - Also replace the url end point.
  - The user name and password of the login dashbord will relaced in "credential.json" file.
