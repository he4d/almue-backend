# github.com/he4d/almue

Welcome to the Almue generated docs.

## Routes

<details>
<summary>`/*`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/***
	- _GET_
		- [almue.fileServer.func1](/almue/almue.go#L235)

</details>
<details>
<summary>`/api/v1/floors`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/floors**
			- **/**
				- _GET_
					- [almue.(*Almue).(github.com/he4d/almue/almue.getAllFloors)-fm](/almue/almue.go#L182)
				- _POST_
					- [almue.(*Almue).(github.com/he4d/almue/almue.createFloor)-fm](/almue/almue.go#L183)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L185)
				- **/**
					- _PUT_
						- [almue.(*Almue).(github.com/he4d/almue/almue.updateFloor)-fm](/almue/almue.go#L187)
					- _DELETE_
						- [almue.(*Almue).(github.com/he4d/almue/almue.deleteFloor)-fm](/almue/almue.go#L188)
					- _GET_
						- [almue.(*Almue).(github.com/he4d/almue/almue.getFloor)-fm](/almue/almue.go#L186)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}/lightings`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L185)
				- **/lightings**
					- **/**
						- _GET_
							- [almue.(*Almue).(github.com/he4d/almue/almue.getAllLightingsOfFloor)-fm](/almue/almue.go#L203)
						- _POST_
							- [almue.(*Almue).(github.com/he4d/almue/almue.createLighting)-fm](/almue/almue.go#L170)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}/lightings/{lightingID:[0-9]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L185)
				- **/lightings**
					- **/{lightingID:[0-9]+$}**
						- [almue.(*Almue).(github.com/he4d/almue/almue.lightingCtx)-fm](/almue/almue.go#L172)
						- **/**
							- _GET_
								- [almue.(*Almue).(github.com/he4d/almue/almue.getLighting)-fm](/almue/almue.go#L173)
							- _PUT_
								- [almue.(*Almue).(github.com/he4d/almue/almue.updateLighting)-fm](/almue/almue.go#L174)
							- _DELETE_
								- [almue.(*Almue).(github.com/he4d/almue/almue.deleteLighting)-fm](/almue/almue.go#L175)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}/lightings/{lightingID:[0-9]+$}/{action:[a-z]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L185)
				- **/lightings**
					- **/{lightingID:[0-9]+$}**
						- [almue.(*Almue).(github.com/he4d/almue/almue.lightingCtx)-fm](/almue/almue.go#L172)
						- **/{action:[a-z]+$}**
							- **/**
								- _POST_
									- [almue.(*Almue).(github.com/he4d/almue/almue.controlLighting)-fm](/almue/almue.go#L177)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}/shutters`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L185)
				- **/shutters**
					- **/**
						- _POST_
							- [almue.(*Almue).(github.com/he4d/almue/almue.createShutter)-fm](/almue/almue.go#L157)
						- _GET_
							- [almue.(*Almue).(github.com/he4d/almue/almue.getAllShuttersOfFloor)-fm](/almue/almue.go#L190)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}/shutters/{shutterID:[0-9]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L185)
				- **/shutters**
					- **/{shutterID:[0-9]+$}**
						- [almue.(*Almue).(github.com/he4d/almue/almue.shutterCtx)-fm](/almue/almue.go#L159)
						- **/**
							- _DELETE_
								- [almue.(*Almue).(github.com/he4d/almue/almue.deleteShutter)-fm](/almue/almue.go#L162)
							- _GET_
								- [almue.(*Almue).(github.com/he4d/almue/almue.getShutter)-fm](/almue/almue.go#L160)
							- _PUT_
								- [almue.(*Almue).(github.com/he4d/almue/almue.updateShutter)-fm](/almue/almue.go#L161)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}/shutters/{shutterID:[0-9]+$}/{action:[a-z]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L185)
				- **/shutters**
					- **/{shutterID:[0-9]+$}**
						- [almue.(*Almue).(github.com/he4d/almue/almue.shutterCtx)-fm](/almue/almue.go#L159)
						- **/{action:[a-z]+$}**
							- **/**
								- _POST_
									- [almue.(*Almue).(github.com/he4d/almue/almue.controlShutter)-fm](/almue/almue.go#L164)

</details>
<details>
<summary>`/api/v1/lightings`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/lightings**
			- **/**
				- _GET_
					- [almue.(*Almue).(github.com/he4d/almue/almue.getAllLightings)-fm](/almue/almue.go#L169)
				- _POST_
					- [almue.(*Almue).(github.com/he4d/almue/almue.createLighting)-fm](/almue/almue.go#L170)

</details>
<details>
<summary>`/api/v1/lightings/{lightingID:[0-9]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/lightings**
			- **/{lightingID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.lightingCtx)-fm](/almue/almue.go#L172)
				- **/**
					- _PUT_
						- [almue.(*Almue).(github.com/he4d/almue/almue.updateLighting)-fm](/almue/almue.go#L174)
					- _DELETE_
						- [almue.(*Almue).(github.com/he4d/almue/almue.deleteLighting)-fm](/almue/almue.go#L175)
					- _GET_
						- [almue.(*Almue).(github.com/he4d/almue/almue.getLighting)-fm](/almue/almue.go#L173)

</details>
<details>
<summary>`/api/v1/lightings/{lightingID:[0-9]+$}/{action:[a-z]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/lightings**
			- **/{lightingID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.lightingCtx)-fm](/almue/almue.go#L172)
				- **/{action:[a-z]+$}**
					- **/**
						- _POST_
							- [almue.(*Almue).(github.com/he4d/almue/almue.controlLighting)-fm](/almue/almue.go#L177)

</details>
<details>
<summary>`/api/v1/manage/db/backup`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/manage**
			- **/db**
				- **/backup**
					- _GET_
						- [almue.(*Almue).(github.com/he4d/almue/almue.retrieveStoreBackup)-fm](/almue/almue.go#L151)
					- _POST_
						- [almue.(*Almue).(github.com/he4d/almue/almue.restoreStoreBackup)-fm](/almue/almue.go#L152)

</details>
<details>
<summary>`/api/v1/manage/logfile`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/manage**
			- **/logfile**
				- _GET_
					- [almue.(*Almue).(github.com/he4d/almue/almue.getLogfile)-fm](/almue/almue.go#L149)

</details>
<details>
<summary>`/api/v1/shutters`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/shutters**
			- **/**
				- _GET_
					- [almue.(*Almue).(github.com/he4d/almue/almue.getAllShutters)-fm](/almue/almue.go#L156)
				- _POST_
					- [almue.(*Almue).(github.com/he4d/almue/almue.createShutter)-fm](/almue/almue.go#L157)

</details>
<details>
<summary>`/api/v1/shutters/{shutterID:[0-9]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/shutters**
			- **/{shutterID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.shutterCtx)-fm](/almue/almue.go#L159)
				- **/**
					- _GET_
						- [almue.(*Almue).(github.com/he4d/almue/almue.getShutter)-fm](/almue/almue.go#L160)
					- _PUT_
						- [almue.(*Almue).(github.com/he4d/almue/almue.updateShutter)-fm](/almue/almue.go#L161)
					- _DELETE_
						- [almue.(*Almue).(github.com/he4d/almue/almue.deleteShutter)-fm](/almue/almue.go#L162)

</details>
<details>
<summary>`/api/v1/shutters/{shutterID:[0-9]+$}/{action:[a-z]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L69)
		- **/shutters**
			- **/{shutterID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.shutterCtx)-fm](/almue/almue.go#L159)
				- **/{action:[a-z]+$}**
					- **/**
						- _POST_
							- [almue.(*Almue).(github.com/he4d/almue/almue.controlShutter)-fm](/almue/almue.go#L164)

</details>

Total # of routes: 17
