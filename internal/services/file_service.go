func (s *FileService) CreateFile(ctx *gin.Context, name, path, url, contentType, ownerID string, size int64) (*models.File, error) {
	if name == "" || path == "" || url == "" || contentType == "" || ownerID == "" {
		merrors.BadRequest(ctx, "Name, path, URL, content type, and owner ID are required")
		return nil, errors.New("name, path, URL, content type, and owner ID are required")
	file := models.NewFile(fileID, name, path, url, contentType, ownerID, size)
func (s *FileService) UpdateFile(ctx *gin.Context, file *models.File, name, path, url, contentType string, size int64) (*models.File, error) {
	if name == "" && path == "" && url == "" && contentType == "" && size == 0 {
	file.UpdateFile(name, path, url, contentType, size)
	err := models.DeleteFileByID(file.ID) // Ensure the actual deletion is handled in the database.
	// In a real-world scenario, you'd also delete the file from storage.