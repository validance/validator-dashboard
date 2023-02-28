# Validator Dashboard

## features
1. Track address status - new, existing, leave, return (labeling has to be done manually)
2. Track delegation amount by each address
3. Generate delegation summary with histories
4. Track grant rewards
5. Track validator income history
6. Record token prices

## Setup
1. Migrate database
    ```shell
    cd app/db/migrations
    sql-migrate up -env="production"
    ```
   
2. Build binary
   ```shell
   go build -o ../dashboard -C app
   ```
   
3. Run dashboard binary
   - set values in config.yaml
   - run the binary at the same directory with config.yaml