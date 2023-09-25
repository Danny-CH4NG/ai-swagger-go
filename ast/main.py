#Note: The openai-python library support for Azure OpenAI is in preview.
import openai
import json
openai.api_type = "azure"
openai.api_base = ""
openai.api_version = "2023-05-15"
openai.api_key = ""

with open('drafts.json', 'r', encoding='utf-8') as file:
    # 讀取json檔案的內容
    data = json.load(file)
data_str = json.dumps(data, ensure_ascii=False, indent=4)

response = openai.ChatCompletion.create(
    engine="gpt-4-8k", 
    messages=[
        {"role": "system", "content": "使用提供的資料，生成一個符合openapi3的yaml文件，並盡量正確的填補空缺的描述資訊，注意Response欄位可作為Example Value。"},
        {"role": "user", "content": data_str},
    ]
)

print(response['choices'][0]['message']['content'])