# github.com/he4d/almue

Welcome to the Almue generated docs.

## Routes

<details>
<summary>`/*`</summary>

- **/***
	- _GET_
		- [(*Mux).FileServer.func1](https://github.com/pressly/chi/mux.go#L317)

</details>
<details>
<summary>`/api/floors`</summary>

- **/api/floors**
	- **/**
		- _POST_
			- [almue.(*Almue).(github.com/he4d/almue/almue.createFloor)-fm](/almue/almue.go#L165)
		- _GET_
			- [almue.(*Almue).(github.com/he4d/almue/almue.getAllFloors)-fm](/almue/almue.go#L164)

</details>
<details>
<summary>`/api/floors/:floorID`</summary>

- **/api/floors**
	- **/:floorID**
		- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L167)
		- **/**
			- _GET_
				- [almue.(*Almue).(github.com/he4d/almue/almue.getFloor)-fm](/almue/almue.go#L168)
			- _PUT_
				- [almue.(*Almue).(github.com/he4d/almue/almue.updateFloor)-fm](/almue/almue.go#L169)
			- _DELETE_
				- [almue.(*Almue).(github.com/he4d/almue/almue.deleteFloor)-fm](/almue/almue.go#L170)

</details>
<details>
<summary>`/api/floors/:floorID/lightings`</summary>

- **/api/floors**
	- **/:floorID**
		- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L167)
		- **/lightings**
			- **/**
				- _POST_
					- [almue.(*Almue).(github.com/he4d/almue/almue.createLighting)-fm](/almue/almue.go#L187)
				- _GET_
					- [almue.(*Almue).(github.com/he4d/almue/almue.getAllLightingsOfFloor)-fm](/almue/almue.go#L186)

</details>
<details>
<summary>`/api/floors/:floorID/lightings/:lightingID`</summary>

- **/api/floors**
	- **/:floorID**
		- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L167)
		- **/lightings**
			- **/:lightingID**
				- [almue.(*Almue).(github.com/he4d/almue/almue.lightingCtx)-fm](/almue/almue.go#L189)
				- **/**
					- _PUT_
						- [almue.(*Almue).(github.com/he4d/almue/almue.updateLighting)-fm](/almue/almue.go#L191)
					- _DELETE_
						- [almue.(*Almue).(github.com/he4d/almue/almue.deleteLighting)-fm](/almue/almue.go#L192)
					- _GET_
						- [almue.(*Almue).(github.com/he4d/almue/almue.getLighting)-fm](/almue/almue.go#L190)

</details>
<details>
<summary>`/api/floors/:floorID/lightings/:lightingID/:action`</summary>

- **/api/floors**
	- **/:floorID**
		- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L167)
		- **/lightings**
			- **/:lightingID**
				- [almue.(*Almue).(github.com/he4d/almue/almue.lightingCtx)-fm](/almue/almue.go#L189)
				- **/:action**
					- [almue.(*Almue).(github.com/he4d/almue/almue.deviceActionCtx)-fm](/almue/almue.go#L180)
					- **/**
						- _POST_
							- [almue.(*Almue).(github.com/he4d/almue/almue.controlLighting)-fm](/almue/almue.go#L195)

</details>
<details>
<summary>`/api/floors/:floorID/shutters`</summary>

- **/api/floors**
	- **/:floorID**
		- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L167)
		- **/shutters**
			- **/**
				- _GET_
					- [almue.(*Almue).(github.com/he4d/almue/almue.getAllShuttersOfFloor)-fm](/almue/almue.go#L172)
				- _POST_
					- [almue.(*Almue).(github.com/he4d/almue/almue.createShutter)-fm](/almue/almue.go#L173)

</details>
<details>
<summary>`/api/floors/:floorID/shutters/:shutterID`</summary>

- **/api/floors**
	- **/:floorID**
		- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L167)
		- **/shutters**
			- **/:shutterID**
				- [almue.(*Almue).(github.com/he4d/almue/almue.shutterCtx)-fm](/almue/almue.go#L175)
				- **/**
					- _GET_
						- [almue.(*Almue).(github.com/he4d/almue/almue.getShutter)-fm](/almue/almue.go#L176)
					- _PUT_
						- [almue.(*Almue).(github.com/he4d/almue/almue.updateShutter)-fm](/almue/almue.go#L177)
					- _DELETE_
						- [almue.(*Almue).(github.com/he4d/almue/almue.deleteShutter)-fm](/almue/almue.go#L178)

</details>
<details>
<summary>`/api/floors/:floorID/shutters/:shutterID/:action`</summary>

- **/api/floors**
	- **/:floorID**
		- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L167)
		- **/shutters**
			- **/:shutterID**
				- [almue.(*Almue).(github.com/he4d/almue/almue.shutterCtx)-fm](/almue/almue.go#L175)
				- **/:action**
					- [almue.(*Almue).(github.com/he4d/almue/almue.deviceActionCtx)-fm](/almue/almue.go#L180)
					- **/**
						- _POST_
							- [almue.(*Almue).(github.com/he4d/almue/almue.controlShutter)-fm](/almue/almue.go#L181)

</details>
<details>
<summary>`/api/lightings`</summary>

- **/api/lightings**
	- **/**
		- _GET_
			- [almue.(*Almue).(github.com/he4d/almue/almue.getAllLightings)-fm](/almue/almue.go#L161)

</details>
<details>
<summary>`/api/shutters`</summary>

- **/api/shutters**
	- **/**
		- _GET_
			- [almue.(*Almue).(github.com/he4d/almue/almue.getAllShutters)-fm](/almue/almue.go#L158)

</details>

Total # of routes: 11
