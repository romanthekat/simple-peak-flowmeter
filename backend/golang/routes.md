# github.com/EvilKhaosKat/simple-peak-flowmeter/backend/golang

Simple peak flowmeter golang generated docs.

## Routes

<details>
<summary>`/ping`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- [URLFormat]()
- [SetContentType.func1]()
- **/ping**
	- _GET_
		- [main.main.func1]()

</details>
<details>
<summary>`/records/*`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- [URLFormat]()
- [SetContentType.func1]()
- **/records/***
	- **/**
		- _POST_
			- [main.CreateRecord]()
		- _GET_
			- [main.ListRecords]()

</details>
<details>
<summary>`/records/*/{RecordID}/*`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- [URLFormat]()
- [SetContentType.func1]()
- **/records/***
	- **/{RecordID}/***
		- [main.RecordCtx]()
		- **/**
			- _PUT_
				- [main.UpdateRecord]()
			- _DELETE_
				- [main.DeleteRecord]()
			- _GET_
				- [main.GetRecord]()

</details>

Total # of routes: 3

