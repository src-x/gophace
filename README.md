```
  ___   __   ____  _  _   __    ___  ____ 
 / __) /  \ (  _ \/ )( \ / _\  / __)(  __)
( (_ \(  O ) ) __/) __ (/    \( (__  ) _) 
 \___/ \__/ (__)  \_)(_/\_/\_/ \___)(____)

```

## Demo
![Demo of gophace](images/demo.gif)

## Technologies
NATS Streaming, GoCV, and Ebiten

## File Structure
```
gophace/
|-- images
|   |-- demo.gif
|   `-- gopher_template.png
|-- publish
|   |-- demo.go
|   `-- publish.go
|-- subscribe
|   `-- subscribe.go
|-- xml_files
|   |-- haarcascade_animeface2.xml
|   |-- haarcascade_frontalface_alt.xml
|   |-- haarcascade_frontalface_default.xml
|   `-- lbpcascade_animeface.xml
```

## Getting Started
```
go run publish/publish.go
```
```
go run subscribe/subscribe.go
```

## Future Improvements and Additions
* Get video stream to work on Gopherjs
* Add face recognition with Spark Streaming