CREATE OR REPLACE FUNCTION public.notify_event()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$

    DECLARE 
        data json;
        notification json;
    
    BEGIN
    
        data = row_to_json(NEW);
       
        notification = json_build_object('data', data);
                        
        -- Execute pg_notify(channel, notification)
        PERFORM pg_notify('events',notification::text);
        
        -- Result is ignored since this is an AFTER trigger
        RETURN NULL; 
    END;
    
$function$;


CREATE TABLE public."data" (
	id serial NOT NULL,
	"date" timestamp NOT NULL DEFAULT now(),
	description varchar(200) NULL
)
WITH (
	OIDS=FALSE
) ;


create
    trigger data_notify_event after insert
        on
        data for each row execute procedure notify_event();
