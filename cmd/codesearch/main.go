
func NewCodeSearch(dic *di.Container) *cobra.Command {
	var dir string
	var pattern string

	cmd := &cobra.Command{
		Use:   "code-search",
		Short: "code-search は指定した path のフォルダを再帰的にコード検索し、指定したパターンにマッチする行を表示します",
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := searchCodeRegex(dir, pattern)
			if err != nil {
				return fmt.Errorf("searching code: %v", err)
			}
			fmt.Print(res)
			return nil
		},
	}
	cmd.Flags().StringVarP(&dir, "directory", "d", "", "Directory to search")
	cmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Regex pattern to search for")
	return cmd
}

func searchCodeRegex(dir, pattern string) (string, error) {
	if dir == "" || pattern == "" {
		return "", fmt.Errorf("both directory and pattern must be specified")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regex pattern: %s %v", re, err)
	}

	sb := strings.Builder{}

	// header
	sb.WriteString("filepath\tline number\tmatch\tline\n")
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("error open file %q: %v", path, err)
			}
			defer file.Close()

			reader := bufio.NewReader(file)
			lineNumber := 1
			for {
				line, err := reader.ReadString('\n')
				if err != nil && err != io.EOF {
					return fmt.Errorf("error reading file %q: %v", path, err)
				}
				if err == io.EOF {
					break
				}
				matches := re.FindAllString(line, -1)
				for _, match := range matches {
					nlrmLine := strings.Trim(line, " \n\t")
					sb.WriteString(fmt.Sprintf("%s\t%d\t%s\t%s\n", path, lineNumber, match, nlrmLine))
				}
				lineNumber++
			}
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("error walking the path %q: %v", dir, err)
	}
	return sb.String(), nil
}