package aggregate

import "contrib/jsternberg/math"

// table will aggregate columns and create tables with a single
// row containing the aggregated value.
//
// This function takes a single parameter of `columns`. The parameter
// is an object with the output column name as the key and the aggregate
// object as the value.
//
// The aggregate object is composed of at least the following required attributes:
//     column = string
//         The column name for the input.
//     init = (values) -> state
//         An initial function to compute the initial state of the
//         output. This can return either the final aggregate or a
//         temporary state object that can be used to compute the
//         final aggregate. The values parameter will always be a
//         non-empty array of values from the specified column.
//     reduce = (values, state) -> state
//         A function that takes in another buffer of values
//         and the current state of the aggregate and computes
//         the updated state.
//     compute = (state) -> value
//         A function that takes the state and computes the final
//         aggregate.
//
// An example of usage is:
//     tables |> aggregate.table(columns: {
//         "min_bottom_degrees": aggregate.min(column: "bottom_degrees"),
//     ])
builtin table

_make_selector = (column, fn) => ({
	column: column,
	init: fn,
	reduce: (values, state) => {
		v = fn(values)
		return fn(values: [state, v])
	},
	compute: (state) => state,
})

// min constructs a min aggregate or selector for the column.
min = (column) => _make_selector(column: column, fn: math.min)

// max constructs a max aggregate or selector for the column.
max = (column) => _make_selector(column: column, fn: math.max)

// sum constructs a sum aggregate for the column.
sum = (column) => ({
	init: (values) => math.sum(values),
	reduce: (values, state) => {
		return state + math.sum(values)
	},
	compute: (state) => state,
})

// count constructs a count aggregate for the column.
count = (column) => ({
	column: column,
	init: (values) => length(arr: values),
	reduce: (values, state) => {
		return state + length(arr: values)
	},
	compute: (state) => state,
})

// mean constructs a mean aggregate for the column.
mean = (column) => ({
	column: column,
	init: (values) => ({
		sum: math.sum(values),
		count: length(arr: values),
	}),
	reduce: (values, state) => ({
		sum: state.sum + math.sum(values),
		count: state.count + length(arr: values),
	}),
	compute: (state) => float(v: state.sum) / float(v: state.count),
})
