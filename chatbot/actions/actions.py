# This files contains your custom actions which can be used to run
# custom Python code.
#
# See this guide on how to implement these action:
# https://rasa.com/docs/rasa/custom-actions


# This is a simple example for a custom action which utters "Hello World!"

from typing import Any, Text, Dict, List
#
from rasa_sdk import Action, Tracker
from rasa_sdk.executor import CollectingDispatcher
import random
import math
import os
# poor man's mysql client
mysql_client = """mysql -uroot -p0000 -s -h127.0.0.1 mysql -e "{}" """
local_tidb_client = """mysql --comments --host 127.0.0.1 --port 4000 -u root mysql -e "{}" """

def run_sql(cmd):
    if os.environ.get('usetidb') != "":
        return "".join(os.popen(local_tidb_client.format(cmd)).readlines())
    return "".join(os.popen(mysql_client.format(cmd)).readlines())
def randomPickGifUrl():
    seed = [
        "https://n.sinaimg.cn/tech/transform/340/w162h178/20210629/aed8-krwipar9961537.gif",
        "https://n.sinaimg.cn/tech/transform/279/w136h143/20210629/bb5f-krwipar9958933.gif",
        "https://n.sinaimg.cn/tech/transform/298/w183h115/20210629/a5e5-krwipar9952310.gif",
        "https://n.sinaimg.cn/tech/transform/396/w207h189/20220105/b668-ddc78350ce0f1db6275482c24490768f.gif",
        "http://localhost:9090/1.png",
        "http://localhost:9090/2.png",
        "http://localhost:9090/3.png",
        "http://localhost:9090/4.png",
        "http://localhost:9090/5.png",
        "http://localhost:9090/6.png",
        "http://localhost:9090/7.png",
        "http://localhost:9090/8.png",
        "http://localhost:9090/9.png",
        "http://localhost:9090/10.png",
    ]
    idx = math.floor(random.random()*10000) % len(seed)
    return seed[idx]

def randomPickGif(dispatcher):
    imgurl = randomPickGifUrl()
    dispatcher.utter_message(image=imgurl)
class ActionHelloWorld(Action):

    def name(self) -> Text:
        return "action_hello_world"

    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        dispatcher.utter_message(text="Hello World!")
        randomPickGif(dispatcher)
        return []
    
class ActionShowDatabases(Action):
    def name(self) -> Text:
        return "action_show_databases"
    
    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        dispatcher.utter_message(text="我有"+"  ".join(run_sql("show databases").lower().split("\n")[1:])+"这几个库呢.")
        randomPickGif(dispatcher)
        return []
    
class ActionShowTables(Action):
    def name(self) -> Text:
        return "action_show_tables"
    
    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        dispatcher.utter_message(text="我有"+" ".join(run_sql("show tables").lower().split("\n")[1:]) + "这些表呢")
        randomPickGif(dispatcher)
        return []
    
class ActionShowStatus(Action):
    def name(self) -> Text:
        return "action_show_status"
    
    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:
        cmd = "status"
        rep = run_sql(cmd)
        dispatcher.utter_message(text=rep)
        randomPickGif(dispatcher)
        return []
    
class ActionExplain(Action):
    def name(self) -> Text:
        return "action_do_explain"
    
    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:
        print(tracker.get_slot("sql"))
        print(tracker.latest_message)
        cmd = "explain {}".format(tracker.get_slot("sql"))
        rep = run_sql(cmd)
        dispatcher.utter_message(text=rep)
        randomPickGif(dispatcher)
        return []


class ActionRunSql(Action):
    def name(self) -> Text:
        return "action_run_sql"
    
    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:
        print(tracker.get_slot("sql"))
        print(tracker.latest_message)
        cmd = "{}".format(tracker.get_slot("sql"))
        rep = run_sql(cmd)
        dispatcher.utter_message(text=rep)
        randomPickGif(dispatcher)
        return []

    
class ActionShowProcesslist(Action):
    def name(self) -> Text:
        return "action_show_processlist"
    
    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        dispatcher.utter_message(text=run_sql("show processlist"))
        randomPickGif(dispatcher)
        return []

class ShowConnection(Action):
    def name(self) -> Text:
        return "action_show_connection"

    def run(self,dispatcher: CollectingDispatcher,tracker: Tracker,
        domain: Dict[Text,Any]
    ) -> List[Dict[Text,Any]]:
        cmd = "select max(value) from METRICS_SCHEMA.tidb_connection_count;"
        rep = run_sql(cmd).lower().split("\n")[1] + "个哦"
        dispatcher.utter_message(text=rep)
        randomPickGif(dispatcher)