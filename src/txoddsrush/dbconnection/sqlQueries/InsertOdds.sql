INSERT INTO public.fact_odds (odds_id, bookie_key, match_key, 
                            order_of_offer, hteam_to_win, ateam_to_win, 
                            draw, odds_last_updated, table_last_updated, bookie_margin,
                            actual_hteam_win, actual_ateam_win, actual_draw) 
                            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)