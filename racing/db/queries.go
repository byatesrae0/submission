package db

const (
	racesList = "list"
)

func getRaceQueries() map[string]string {
	return map[string]string{
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
}
