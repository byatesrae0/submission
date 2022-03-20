package db

const (
	racesList = "list"
	raceGet   = "get"
)

func getRaceQueries() map[string]string {
	m := map[string]string{
		racesList: `
			SELECT 
				id, 
				meeting_id, 
				name, 
				number, 
				visible, 
				advertised_start_time,
				CASE WHEN strftime('%s', advertised_start_time) >= strftime('%s','now')
             	THEN 1 ELSE 2 END AS status
			FROM races
		`,
	}

	m[raceGet] = m[racesList]

	return m
}
