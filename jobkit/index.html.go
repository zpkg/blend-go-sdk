package jobkit

var indexTemplate = `
{{ define "index" }}
{{ template "header" . }}
<div class="container">
	<table class="u-full-width">
		<thead>
			<tr>
				<th>Job Name</th>
				<th>Status</th>
				<th>Current</th>
				<th>Next Run</th>
				<th>Last Ran</th>
				<th>Last Result</th>
				<th>Last Elapsed</th>
			</tr>
		</thead>
		<tbody>
		{{ range $index, $job := .ViewModel.Jobs }}
			<tr>
				<td> <!-- job name -->
					{{ $job.Name }}
				</td>
				<td> <!-- job status -->
				{{ if $job.Disabled }}
					<span class="danger">Disabled</span>
				{{else}}
					<span class="primary">Enabled</span>
				{{end}}
				</td>
				<td> <!-- job status -->
				{{ if $job.Current }}
					{{ since_utc $job.Current.Started }}
				{{else}}
					<span>-</span>
				{{end}}
				</td>
				<td> <!-- next run-->
					{{ $job.NextRuntime | rfc3339 }}
				</td>
				<td> <!-- last run -->
				{{ if $job.Last }}
					{{ $job.Last.Started | rfc3339 }}
				{{ else }}
					<span class="none">-</span>
				{{ end }}
				</td>
				<td> <!-- last status -->
				{{ if $job.Last }}
					{{ if $job.Last.Err }}
						{{ $job.Last.Err }}
					{{ else }}
					<span class="none">Success</span>
					{{ end }}
				{{ else }}
					<span class="none">-</span>
				{{ end }}
				</td>
				<td><!-- last elapsed -->
				{{ if $job.Last }}
					{{ $job.Last.Elapsed }}
				{{ else }}
					<span class="none">-</span>
				{{ end }}
				</td>
			</tr>
			<tr>
				<td colspan=6>
					<table class="u-full-width">
						{{ range $index, $ji := $job.History }}
						<tr>
							<td>{{ $ji.ID }}</td>
							<td>{{ $ji.Started | rfc3339 }}</td>
							<td>{{ $ji.Finished | rfc3339 }}</td>
							<td>{{ $ji.Timeout | rfc3339 }}</td>
							<td>{{ $ji.Cancelled | rfc3339 }}</td>
							<td>{{ $ji.Elapsed }}</td>
							<td>{{ $ji.Err }}</td>
						</tr>
						{{ end }}
					</table>
				</td>
			</tr>
		{{ else }}
			<tr><td colspan=6>No Jobs Loaded</td></tr>
		{{ end }}
		</tbody>
	</table>
</div>
{{ template "footer" . }}
{{ end }}
`
