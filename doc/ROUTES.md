# github.com/he4d/almue

Welcome to the Almue generated docs.

## Routes

<details>
<summary>`/*`</summary>

- [RequestID](https://github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/pressly/chi/middleware/recoverer.go#L16)
- **/***
	- _GET_
		- [(*Mux).FileServer.func1](https://github.com/pressly/chi/mux.go#L317)

</details>
<details>
<summary>`/api/v1/floors`</summary>

- [RequestID](https://github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L72)
		- **/floors**
			- **/**
				- _GET_
					- [almue.(*Almue).(github.com/he4d/almue/almue.getAllFloors)-fm](/almue/almue.go#L173)
				- _POST_
					- [almue.(*Almue).(github.com/he4d/almue/almue.createFloor)-fm](/almue/almue.go#L174)

</details>
<details>
<summary>`/api/v1/floors/:floorID`</summary>

- [RequestID](https://github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L72)
		- **/floors**
			- **/:floorID**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L176)
				- **/**
					- _DELETE_
						- [almue.(*Almue).(github.com/he4d/almue/almue.deleteFloor)-fm](/almue/almue.go#L179)
					- _GET_
						- [almue.(*Almue).(github.com/he4d/almue/almue.getFloor)-fm](/almue/almue.go#L177)
					- _PUT_
						- [almue.(*Almue).(github.com/he4d/almue/almue.updateFloor)-fm](/almue/almue.go#L178)

</details>
<details>
<summary>`/api/v1/floors/:floorID/lightings`</summary>

- [RequestID](https://github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L72)
		- **/floors**
			- **/:floorID**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L176)
				- **/lightings**
					- **/**
						- _GET_
							- [almue.(*Almue).(github.com/he4d/almue/almue.getAllLightingsOfFloor)-fm](/almue/almue.go#L194)
						- _POST_
							- [almue.(*Almue).(github.com/he4d/almue/almue.createLighting)-fm](/almue/almue.go#L195)

</details>
<details>
<summary>`/api/v1/floors/:floorID/lightings/:lightingID`</summary>

- [RequestID](https://github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L72)
		- **/floors**
			- **/:floorID**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L176)
				- **/lightings**
					- **/:lightingID**
						- [almue.(*Almue).(github.com/he4d/almue/almue.lightingCtx)-fm](/almue/almue.go#L197)
						- **/**
							- _GET_
								- [almue.(*Almue).(github.com/he4d/almue/almue.getLighting)-fm](/almue/almue.go#L198)
							- _PUT_
								- [almue.(*Almue).(github.com/he4d/almue/almue.updateLighting)-fm](/almue/almue.go#L199)
							- _DELETE_
								- [almue.(*Almue).(github.com/he4d/almue/almue.deleteLighting)-fm](/almue/almue.go#L200)

</details>
<details>
<summary>`/api/v1/floors/:floorID/lightings/:lightingID/:action`</summary>

- [RequestID](https://github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L72)
		- **/floors**
			- **/:floorID**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L176)
				- **/lightings**
					- **/:lightingID**
						- [almue.(*Almue).(github.com/he4d/almue/almue.lightingCtx)-fm](/almue/almue.go#L197)
						- **/:action**
							- **/**
								- _POST_
									- [almue.(*Almue).(github.com/he4d/almue/almue.controlLighting)-fm](/almue/almue.go#L202)

</details>
<details>
<summary>`/api/v1/floors/:floorID/shutters`</summary>

- [RequestID](https://github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L72)
		- **/floors**
			- **/:floorID**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L176)
				- **/shutters**
					- **/**
						- _GET_
							- [almue.(*Almue).(github.com/he4d/almue/almue.getAllShuttersOfFloor)-fm](/almue/almue.go#L181)
						- _POST_
							- [almue.(*Almue).(github.com/he4d/almue/almue.createShutter)-fm](/almue/almue.go#L182)

</details>
<details>
<summary>`/api/v1/floors/:floorID/shutters/:shutterID`</summary>

- [RequestID](https://github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L72)
		- **/floors**
			- **/:floorID**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L176)
				- **/shutters**
					- **/:shutterID**
						- [almue.(*Almue).(github.com/he4d/almue/almue.shutterCtx)-fm](/almue/almue.go#L184)
						- **/**
							- _GET_
								- [almue.(*Almue).(github.com/he4d/almue/almue.getShutter)-fm](/almue/almue.go#L185)
							- _PUT_
								- [almue.(*Almue).(github.com/he4d/almue/almue.updateShutter)-fm](/almue/almue.go#L186)
							- _DELETE_
								- [almue.(*Almue).(github.com/he4d/almue/almue.deleteShutter)-fm](/almue/almue.go#L187)

</details>
<details>
<summary>`/api/v1/floors/:floorID/shutters/:shutterID/:action`</summary>

- [RequestID](https://github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L72)
		- **/floors**
			- **/:floorID**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L176)
				- **/shutters**
					- **/:shutterID**
						- [almue.(*Almue).(github.com/he4d/almue/almue.shutterCtx)-fm](/almue/almue.go#L184)
						- **/:action**
							- **/**
								- _POST_
									- [almue.(*Almue).(github.com/he4d/almue/almue.controlShutter)-fm](/almue/almue.go#L189)

</details>
<details>
<summary>`/api/v1/lightings`</summary>

- [RequestID](https://github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L72)
		- **/lightings**
			- **/**
				- _GET_
					- [almue.(*Almue).(github.com/he4d/almue/almue.getAllLightings)-fm](/almue/almue.go#L170)

</details>
<details>
<summary>`/api/v1/shutters`</summary>

- [RequestID](https://github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L72)
		- **/shutters**
			- **/**
				- _GET_
					- [almue.(*Almue).(github.com/he4d/almue/almue.getAllShutters)-fm](/almue/almue.go#L167)

</details>

Total # of routes: 11
