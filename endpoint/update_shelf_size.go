package endpoint

/*
func CreateUpdateShelfSizeEndpoint(
	updateShelfSize shelf.UpdateShelfSizeFunc,
) UpdateShelfSizeEndpoint {
	return func(ctx context.Context, req *gateway.UpdateShelfSizeRequest) error {
		handleError := func(err error) error {
			return fmt.Errorf("on update shelf size endpoint: %w", err)
		}
		userId := core.UserId(req.UserId)
		size := shelf.Size(req.Size)
		err := updateShelfSizeRepo(ctx, userId, size)
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
*/
