# PromQL

- [Instant Vectors](#instant-vectors)
- [Range Vectors](#range-vectors)
- [Operators](#operators)

Note: [PromLens](https://demo.promlens.com/) is a good place to try out different queries

## Instant Vectors

> a set of time series containing a single sample for each time series, all sharing the same timestamp.

```promql
api_invoke_count
api_invoke_count{api_name="beer"}
api_invoke_count{api_name="beer",success="true"}
api_invoke_count{api_name!="beer",success="true"}
api_invoke_count{api_name=~"car|delay"}
```

## Range Vectors

> a set of time series containing a range of data points over time for each time series, can't be charted.

```promql
api_invoke_count[1m]
api_invoke_count{api_name="beer"}[1m]
api_invoke_count{api_name="beer",success="true"}[1m]
api_invoke_count{api_name!="beer",success="true"}[1m]
api_invoke_count{api_name=~"car|delay"}[1m]
```

## Operators

- `+, -, *, /, % and ^(power)` arithmetic operators (defined between scalar/scaler, vector/scalr and vector/vector).
- `==, >=, <, <=, >, !=` logical operators (defined between scalar/scaler, vector/scalr and vector/vector)..
- `sum`, `max`, `min`, `count`, and `avg` etc... aggregate operators (single instant vector and returns new vector).

> Getting the API Success percentage

```promql
sum(api_invoke_count{success="true"}) * 100 / (sum(api_invoke_count{success="true"}) + sum(api_invoke_count{success="false"}))
```

> Success percentage for beer suggestions >= 70%

```promql
sum(api_invoke_count{success="true", api_name="beer"} * 100) / sum(api_invoke_count{api_name="beer"}) > 70
```

> Getting individual success percentage of all the APIs

```promql
(sum by (api_name) (api_invoke_count{success="true"})) * 100 / ((sum by (api_name) (api_invoke_count{success="false"})) + (sum by (api_name) (api_invoke_count{success="true"})))
```

## References

- [PromQL Basics](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [PromQL cheat sheet](https://promlabs.com/promql-cheat-sheet/)
