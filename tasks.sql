

CREATE UNIQUE INDEX one_accepted_application_per_task
ON task_application(task_id)
WHERE status = 'accepted';

