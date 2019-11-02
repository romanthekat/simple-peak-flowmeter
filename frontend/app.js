var chart_width = 800
var chart_height = 400
var padding = 50


d3.json('http://romangaranin.dev:3333/records').then(function(data) {
	console.log(data)
	generate(data)
})

function generate(data) {
	var svg = d3.select('#chart')
		.append('svg')
		.attr('width', chart_width)
		.attr('height', chart_height)
		.attr('fill', 'none')
    .attr('stroke-linejoin', 'round')
    .attr('stroke-linecap', 'round')

	// convert ISO 8601 time into d3 time object
	data.forEach(function(e, i) {
		data[i].created_at = d3.isoParse(e.created_at)
	})


	// make scales
	var x_scale = d3.scaleTime()
		.domain([
			d3.min(data, d => d.created_at), 
			d3.max(data, d => d.created_at)
		])
		.range([padding, chart_width - padding])

	var y_scale = d3.scaleLinear()
		.domain([
			200,
			d3.max(data, d => d.value)
		])
		.range([chart_height - padding, padding])


	// add axis
	var x_axis = d3.axisBottom(x_scale)
	svg.append('g')
		.attr('class', 'x-axis')
		.attr(
			'transform',
			'translate(0, ' + (chart_height - padding) + ')'
		)
		.call(x_axis)

	var y_axis = d3.axisLeft(y_scale)
	svg.append('g')
		.attr('class', 'y-axis')
		.attr(
			'transform',
			'translate(' + padding + ', 0)'
		)
		.call(y_axis)


	// create line
	var line = d3.line()
    .defined(d => !isNaN(d.value))
    .x(d => x_scale(d.created_at))
    .y(d => y_scale(d.value))

  svg.append('path')
	  .datum(data)
	  .attr('fill', 'none')
	  .attr('stroke', 'steelblue')
	  .attr('stroke-width', 1.5)
	  .attr("stroke-linejoin', 'round')
	  .attr('stroke-linecap', 'round')
	  .attr('d', line);

	// svg.append("path")
	//   .datum(data.filter(line.defined()))
	//   .attr("stroke", "#ccc")
	//   .attr("d", line);
}