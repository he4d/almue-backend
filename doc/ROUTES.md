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
		- [almue.fileServer.func1](/almue/almue.go#L226)

</details>
<details>
<summary>`/api/v1/floors`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/floors**
			- **/**
				- _GET_
					- [almue.(*Almue).(github.com/he4d/almue/almue.getAllFloors)-fm](/almue/almue.go#L174)
				- _POST_
					- [almue.(*Almue).(github.com/he4d/almue/almue.createFloor)-fm](/almue/almue.go#L175)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L177)
				- **/**
					- _DELETE_
						- [almue.(*Almue).(github.com/he4d/almue/almue.deleteFloor)-fm](/almue/almue.go#L180)
					- _GET_
						- [almue.(*Almue).(github.com/he4d/almue/almue.getFloor)-fm](/almue/almue.go#L178)
					- _PUT_
						- [almue.(*Almue).(github.com/he4d/almue/almue.updateFloor)-fm](/almue/almue.go#L179)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}/lightings`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L177)
				- **/lightings**
					- **/**
						- _POST_
							- [almue.(*Almue).(github.com/he4d/almue/almue.createLighting)-fm](/almue/almue.go#L162)
						- _GET_
							- [almue.(*Almue).(github.com/he4d/almue/almue.getAllLightingsOfFloor)-fm](/almue/almue.go#L195)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}/lightings/{lightingID:[0-9]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L177)
				- **/lightings**
					- **/{lightingID:[0-9]+$}**
						- [almue.(*Almue).(github.com/he4d/almue/almue.lightingCtx)-fm](/almue/almue.go#L164)
						- **/**
							- _PUT_
								- [almue.(*Almue).(github.com/he4d/almue/almue.updateLighting)-fm](/almue/almue.go#L166)
							- _DELETE_
								- [almue.(*Almue).(github.com/he4d/almue/almue.deleteLighting)-fm](/almue/almue.go#L167)
							- _GET_
								- [almue.(*Almue).(github.com/he4d/almue/almue.getLighting)-fm](/almue/almue.go#L165)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}/lightings/{lightingID:[0-9]+$}/{action:[a-z]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L177)
				- **/lightings**
					- **/{lightingID:[0-9]+$}**
						- [almue.(*Almue).(github.com/he4d/almue/almue.lightingCtx)-fm](/almue/almue.go#L164)
						- **/{action:[a-z]+$}**
							- **/**
								- _POST_
									- [almue.(*Almue).(github.com/he4d/almue/almue.controlLighting)-fm](/almue/almue.go#L169)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}/shutters`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L177)
				- **/shutters**
					- **/**
						- _GET_
							- [almue.(*Almue).(github.com/he4d/almue/almue.getAllShuttersOfFloor)-fm](/almue/almue.go#L182)
						- _POST_
							- [almue.(*Almue).(github.com/he4d/almue/almue.createShutter)-fm](/almue/almue.go#L149)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}/shutters/{shutterID:[0-9]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L177)
				- **/shutters**
					- **/{shutterID:[0-9]+$}**
						- [almue.(*Almue).(github.com/he4d/almue/almue.shutterCtx)-fm](/almue/almue.go#L151)
						- **/**
							- _GET_
								- [almue.(*Almue).(github.com/he4d/almue/almue.getShutter)-fm](/almue/almue.go#L152)
							- _PUT_
								- [almue.(*Almue).(github.com/he4d/almue/almue.updateShutter)-fm](/almue/almue.go#L153)
							- _DELETE_
								- [almue.(*Almue).(github.com/he4d/almue/almue.deleteShutter)-fm](/almue/almue.go#L154)

</details>
<details>
<summary>`/api/v1/floors/{floorID:[0-9]+$}/shutters/{shutterID:[0-9]+$}/{action:[a-z]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/floors**
			- **/{floorID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.floorCtx)-fm](/almue/almue.go#L177)
				- **/shutters**
					- **/{shutterID:[0-9]+$}**
						- [almue.(*Almue).(github.com/he4d/almue/almue.shutterCtx)-fm](/almue/almue.go#L151)
						- **/{action:[a-z]+$}**
							- **/**
								- _POST_
									- [almue.(*Almue).(github.com/he4d/almue/almue.controlShutter)-fm](/almue/almue.go#L156)

</details>
<details>
<summary>`/api/v1/lightings`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/lightings**
			- **/**
				- _POST_
					- [almue.(*Almue).(github.com/he4d/almue/almue.createLighting)-fm](/almue/almue.go#L162)
				- _GET_
					- [almue.(*Almue).(github.com/he4d/almue/almue.getAllLightings)-fm](/almue/almue.go#L161)

</details>
<details>
<summary>`/api/v1/lightings/{lightingID:[0-9]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/lightings**
			- **/{lightingID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.lightingCtx)-fm](/almue/almue.go#L164)
				- **/**
					- _GET_
						- [almue.(*Almue).(github.com/he4d/almue/almue.getLighting)-fm](/almue/almue.go#L165)
					- _PUT_
						- [almue.(*Almue).(github.com/he4d/almue/almue.updateLighting)-fm](/almue/almue.go#L166)
					- _DELETE_
						- [almue.(*Almue).(github.com/he4d/almue/almue.deleteLighting)-fm](/almue/almue.go#L167)

</details>
<details>
<summary>`/api/v1/lightings/{lightingID:[0-9]+$}/{action:[a-z]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/lightings**
			- **/{lightingID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.lightingCtx)-fm](/almue/almue.go#L164)
				- **/{action:[a-z]+$}**
					- **/**
						- _POST_
							- [almue.(*Almue).(github.com/he4d/almue/almue.controlLighting)-fm](/almue/almue.go#L169)

</details>
<details>
<summary>`/api/v1/shutters`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/shutters**
			- **/**
				- _GET_
					- [almue.(*Almue).(github.com/he4d/almue/almue.getAllShutters)-fm](/almue/almue.go#L148)
				- _POST_
					- [almue.(*Almue).(github.com/he4d/almue/almue.createShutter)-fm](/almue/almue.go#L149)

</details>
<details>
<summary>`/api/v1/shutters/{shutterID:[0-9]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/shutters**
			- **/{shutterID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.shutterCtx)-fm](/almue/almue.go#L151)
				- **/**
					- _PUT_
						- [almue.(*Almue).(github.com/he4d/almue/almue.updateShutter)-fm](/almue/almue.go#L153)
					- _DELETE_
						- [almue.(*Almue).(github.com/he4d/almue/almue.deleteShutter)-fm](/almue/almue.go#L154)
					- _GET_
						- [almue.(*Almue).(github.com/he4d/almue/almue.getShutter)-fm](/almue/almue.go#L152)

</details>
<details>
<summary>`/api/v1/shutters/{shutterID:[0-9]+$}/{action:[a-z]+$}`</summary>

- [RequestID](https://github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](https://github.com/go-chi/chi/middleware/logger.go#L26)
- [Recoverer](https://github.com/go-chi/chi/middleware/recoverer.go#L18)
- [SetContentType.func1](https://github.com/go-chi/chi/render/content_type.go#L49)
- **/api**
	- **/v1**
		- [almue.apiVersionCtx.func1](/almue/context.go#L66)
		- **/shutters**
			- **/{shutterID:[0-9]+$}**
				- [almue.(*Almue).(github.com/he4d/almue/almue.shutterCtx)-fm](/almue/almue.go#L151)
				- **/{action:[a-z]+$}**
					- **/**
						- _POST_
							- [almue.(*Almue).(github.com/he4d/almue/almue.controlShutter)-fm](/almue/almue.go#L156)

</details>

Total # of routes: 15
