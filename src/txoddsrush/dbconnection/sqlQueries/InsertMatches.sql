INSERT INTO public.fact_matches (match_key, group_name, 
                                home_team, away_team, 
                                match_date, match_time, 
                                last_updated, match_timestamp, 
                                table_last_updated) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);