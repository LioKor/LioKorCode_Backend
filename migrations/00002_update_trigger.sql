DROP TRIGGER update_solution ON tasks;

CREATE TRIGGER update_solution
    AFTER UPDATE ON tasks
    FOR EACH ROW
    WHEN (OLD.tests is DISTINCT FROM NEW.tests)
    EXECUTE FUNCTION update_solution();