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
import os
# poor man's mysql client
mysql_client = """mysql -uroot -p0000 -s -h127.0.0.1 mysql -e "{}" """
local_tidb_client = """mysql --comments --host 127.0.0.1 --port 4000 -u root mysql -e "{}" """

def run_sql(cmd):
    if os.environ.get('usetidb') != "":
        return "".join(os.popen(local_tidb_client.format(cmd)).readlines())
    return "".join(os.popen(mysql_client.format(cmd)).readlines())

class ActionHelloWorld(Action):

    def name(self) -> Text:
        return "action_hello_world"

    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        dispatcher.utter_message(text="Hello World!")

        return []
    
class ActionShowDatabases(Action):
    def name(self) -> Text:
        return "action_show_databases"
    
    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        dispatcher.utter_message(text=run_sql("show databases"))

        return []
    
class ActionShowTables(Action):
    def name(self) -> Text:
        return "action_show_tables"
    
    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        dispatcher.utter_message(text=run_sql("show tables"))

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
        return []

    
class ActionShowProcesslist(Action):
    def name(self) -> Text:
        return "action_show_processlist"
    
    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        dispatcher.utter_message(text=run_sql("show processlist"))
        return []

