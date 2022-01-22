INFO	2019/11/03 14:12:09 Connecting to MongoDB
# github.com/romanthekat/simple-peak-flowmeter/backend/golang

Simple peak flowmeter golang generated docs.

## Routes

<details>
<summary>`/records/*`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- [URLFormat]()
- [SetContentType.func1]()
- [(*Cors).Handler-fm]()
- **/records/***
	- **/**
		- _POST_
			- [main.(*application).CreateRecord-fm]()
		- _GET_
			- [main.(*application).ListRecords-fm]()

</details>
<details>
<summary>`/records/*/{RecordID}/*`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- [URLFormat]()
- [SetContentType.func1]()
- [(*Cors).Handler-fm]()
- **/records/***
	- **/{RecordID}/***
		- [main.(*application).RecordCtx-fm]()
		- **/**
			- _PUT_
				- [main.(*application).UpdateRecord-fm]()
			- _DELETE_
				- [main.(*application).DeleteRecord-fm]()
			- _GET_
				- [main.(*application).GetRecord-fm]()

</details>
<details>
<summary>`/static/`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- [URLFormat]()
- [SetContentType.func1]()
- [(*Cors).Handler-fm]()
- **/static/**
	- _*_
		- [StripPrefix.func1]()

</details>

Total # of routes: 3
