import random
from datetime import datetime, timedelta
import ipaddress

def generate_ipv4():
    return str(ipaddress.IPv4Address(random.randint(0, 2**32 - 1)))

def generate_ipv6():
    return str(ipaddress.IPv6Address(random.randint(0, 2**128 - 1)))

def generate_log_entry():
    # Gerar data (últimos 30 dias)
    log_date = datetime.now() - timedelta(days=random.randint(0, 30))
    log_date_str = log_date.strftime("%Y-%m-%d %H:%M:%S")
    
    # Escolher aleatoriamente entre IPv4 e IPv6
    if random.choice([True, False]):
        ip = generate_ipv4()
    else:
        ip = generate_ipv6()
    
    # Grupo do log
    log_group = random.choice(["INFO", "WARNING", "DANGER"])
    
    # Mensagens de exemplo para cada grupo
    messages = {
        "INFO": [
            "Sistema inicializado com sucesso",
            "Backup concluído",
            "Novo usuário conectado",
            "Operação concluída sem erros",
            "Serviço reiniciado conforme agendamento"
        ],
        "WARNING": [
            "Uso de CPU acima de 80%",
            "Tentativa de login falha",
            "Espaço em disco abaixo de 20%",
            "Latência de rede acima do esperado",
            "Temperatura do servidor elevada"
        ],
        "DANGER": [
            "Falha crítica no sistema",
            "Ataque DDoS detectado",
            "Serviço principal indisponível",
            "Violação de segurança detectada",
            "Perda de dados ocorrida"
        ]
    }
    
    message = random.choice(messages[log_group])
    
    return f"{log_date_str}{ip}{log_group}{message}\n"

# Gerar arquivo de log com 100 mil entradas
with open("./tmp/fake_logs.txt", "w") as log_file:
    for _ in range(100_000):
        log_file.write(generate_log_entry())

print("Arquivo de logs fake gerado com sucesso: fake_logs.txt")